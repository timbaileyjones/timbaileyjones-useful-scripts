package db

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	_ "modernc.org/sqlite"
)

// metaRow is the decoded meta key "0" (hex-encoded JSON).
type metaRow struct {
	AgentID          string `json:"agentId"`
	LatestRootBlobID string `json:"latestRootBlobId"`
	Name             string `json:"name"`
	Mode             string `json:"mode"`
	CreatedAt        int64  `json:"createdAt"`
	LastUsedModel    string `json:"lastUsedModel"`
}

// chatBlobMessage is a JSON blob with role/content (chat message).
type chatBlobMessage struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

// contentPart is used when content is an array of parts (e.g. {"type":"text","text":"..."}).
type contentPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// DumpOptions controls dump behavior (e.g. colorization).
type DumpOptions struct {
	Color bool
}

// dbEntry holds a path and its createdAt for sorting.
type dbEntry struct {
	Path      string
	CreatedAt int64
}

// DumpAll discovers all *.db under chatsDir, sorts by createdAt, and dumps each.
func DumpAll(chatsDir string, out io.Writer, opts *DumpOptions) error {
	if opts == nil {
		opts = &DumpOptions{}
	}
	entries, err := discoverDBs(chatsDir)
	if err != nil {
		return fmt.Errorf("discover: %w", err)
	}
	if len(entries) == 0 {
		return nil
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].CreatedAt < entries[j].CreatedAt })
	for _, e := range entries {
		if err := dumpOne(e.Path, out, opts); err != nil {
			return fmt.Errorf("dump %s: %w", e.Path, err)
		}
	}
	return nil
}

// discoverDBs walks chatsDir recursively and returns *.db paths with createdAt (or mtime fallback).
func discoverDBs(chatsDir string) ([]dbEntry, error) {
	var paths []string
	err := filepath.Walk(chatsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".db") && !strings.HasSuffix(path, "-shm") && !strings.HasSuffix(path, "-wal") {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	entries := make([]dbEntry, 0, len(paths))
	for _, p := range paths {
		createdAt := getCreatedAt(p)
		entries = append(entries, dbEntry{Path: p, CreatedAt: createdAt})
	}
	return entries, nil
}

func getCreatedAt(dbPath string) int64 {
	conn, err := openReadOnly(dbPath)
	if err != nil {
		return mtimeMs(dbPath)
	}
	defer conn.Close()
	var value string
	err = conn.QueryRow("SELECT value FROM meta WHERE key = '0'").Scan(&value)
	if err != nil {
		return mtimeMs(dbPath)
	}
	raw, err := hex.DecodeString(value)
	if err != nil {
		return mtimeMs(dbPath)
	}
	var m metaRow
	if json.Unmarshal(raw, &m) != nil {
		return mtimeMs(dbPath)
	}
	return m.CreatedAt
}

func mtimeMs(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.ModTime().UnixMilli()
}

func openReadOnly(dbPath string) (*sql.DB, error) {
	abs, err := filepath.Abs(dbPath)
	if err != nil {
		abs = dbPath
	}
	uri := "file:" + abs + "?mode=ro"
	return sql.Open("sqlite", uri)
}

func dumpOne(dbPath string, out io.Writer, opts *DumpOptions) error {
	conn, err := openReadOnly(dbPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Fprintf(out, "=== %s ===\n", dbPath)

	// Meta table
	rows, err := conn.Query("SELECT key, value FROM meta")
	if err != nil {
		return fmt.Errorf("meta: %w", err)
	}
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			rows.Close()
			return err
		}
		if key == "0" {
			raw, err := hex.DecodeString(value)
			if err != nil {
				fmt.Fprintf(out, "[meta key=%s hex-decode error]\n", key)
				continue
			}
			var m metaRow
			if json.Unmarshal(raw, &m) != nil {
				fmt.Fprintf(out, "[meta key=%s json error]\n", key)
				continue
			}
			ts := time.UnixMilli(m.CreatedAt).Format(time.RFC3339)
			fmt.Fprintf(out, "meta: name=%q agentId=%s createdAt=%s mode=%s lastUsedModel=%s\n",
				m.Name, m.AgentID, ts, m.Mode, m.LastUsedModel)
		} else {
			fmt.Fprintf(out, "[meta key=%s value_len=%d]\n", key, len(value))
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	// Blobs table
	blobRows, err := conn.Query("SELECT id, data FROM blobs")
	if err != nil {
		return fmt.Errorf("blobs: %w", err)
	}
	defer blobRows.Close()
	for blobRows.Next() {
		var id string
		var data []byte
		if err := blobRows.Scan(&id, &data); err != nil {
			return err
		}
		lines := formatBlob(id, data, opts)
		for _, line := range lines {
			fmt.Fprintln(out, line)
		}
	}
	return blobRows.Err()
}

const maxPreviewLen = 800
const maxExtractedLineLen = 400

func formatBlob(id string, data []byte, opts *DumpOptions) []string {
	color := opts != nil && opts.Color
	roleColor := func(role string) (string, string) {
		if !color {
			return "", ""
		}
		switch role {
		case "user":
			return "\033[36m", "\033[0m" // cyan
		case "assistant":
			return "\033[32m", "\033[0m" // green
		case "system":
			return "\033[90m", "\033[0m" // dim
		default:
			return "\033[33m", "\033[0m" // yellow
		}
	}

	if len(data) == 0 {
		return []string{fmt.Sprintf("[blob id=%s len=0]", id)}
	}

	// Binary blob: try to extract UTF-8 strings (protobuf-like).
	if !utf8.Valid(data) {
		extracted := ExtractStringsFromBinary(data)
		if len(extracted) == 0 {
			return []string{fmt.Sprintf("[binary blob id=%s len=%d]", id, len(data))}
		}
		lines := make([]string, 0, 1+len(extracted))
		lines = append(lines, fmt.Sprintf("[binary blob id=%s len=%d] (extracted text below)", id, len(data)))
		for _, s := range extracted {
			if len(s) > maxExtractedLineLen {
				s = s[:maxExtractedLineLen] + "..."
			}
			lines = append(lines, "  │ "+s)
		}
		return lines
	}

	// JSON blob
	var msg chatBlobMessage
	if json.Unmarshal(data, &msg) != nil {
		return []string{fmt.Sprintf("[blob id=%s len=%d]", id, len(data))}
	}

	contentStr := extractContentString(msg.Content)
	if msg.Role != "" && contentStr != "" {
		preview := contentStr
		if len(preview) > maxPreviewLen {
			preview = preview[:maxPreviewLen] + "..."
		}
		preview = strings.TrimSpace(preview)
		open, close := roleColor(msg.Role)
		line := open + msg.Role + close + ": " + preview
		return []string{line}
	}

	// JSON with other shape: show keys and content preview
	if len(msg.Content) > 0 {
		preview := string(msg.Content)
		if len(preview) > 300 {
			preview = preview[:300] + "..."
		}
		return []string{
			fmt.Sprintf("[blob id=%s len=%d] (JSON, role=%q)", id, len(data), msg.Role),
			"  │ " + preview,
		}
	}
	return []string{fmt.Sprintf("[blob id=%s len=%d]", id, len(data))}
}

// extractContentString returns display text from content (string or array of parts).
func extractContentString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	// Try as string first
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	// Try as array of parts
	var parts []contentPart
	if json.Unmarshal(raw, &parts) != nil {
		return ""
	}
	var b strings.Builder
	for _, p := range parts {
		if p.Text != "" {
			if b.Len() > 0 {
				b.WriteString("\n")
			}
			b.WriteString(p.Text)
		}
	}
	return b.String()
}
