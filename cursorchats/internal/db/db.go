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
	Role    string `json:"role"`
	Content string `json:"content"`
}

// dbEntry holds a path and its createdAt for sorting.
type dbEntry struct {
	Path      string
	CreatedAt int64
}

// DumpAll discovers all *.db under chatsDir, sorts by createdAt, and dumps each.
func DumpAll(chatsDir string, out io.Writer) error {
	entries, err := discoverDBs(chatsDir)
	if err != nil {
		return fmt.Errorf("discover: %w", err)
	}
	if len(entries) == 0 {
		return nil
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].CreatedAt < entries[j].CreatedAt })
	for _, e := range entries {
		if err := dumpOne(e.Path, out); err != nil {
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

func dumpOne(dbPath string, out io.Writer) error {
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
		line := formatBlob(id, data)
		fmt.Fprintln(out, line)
	}
	return blobRows.Err()
}

func formatBlob(id string, data []byte) string {
	if len(data) == 0 {
		return fmt.Sprintf("[blob id=%s len=0]", id)
	}
	if !utf8.Valid(data) {
		return fmt.Sprintf("[binary blob id=%s len=%d]", id, len(data))
	}
	var msg chatBlobMessage
	if json.Unmarshal(data, &msg) != nil {
		return fmt.Sprintf("[blob id=%s len=%d]", id, len(data))
	}
	if msg.Role != "" && msg.Content != "" {
		preview := msg.Content
		if len(preview) > 500 {
			preview = preview[:500] + "..."
		}
		return fmt.Sprintf("%s: %s", msg.Role, strings.TrimSpace(preview))
	}
	return fmt.Sprintf("[blob id=%s len=%d]", id, len(data))
}
