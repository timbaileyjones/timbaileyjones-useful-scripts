package db

import (
	"encoding/binary"
	"unicode/utf8"
)

// ExtractStringsFromBinary heuristically extracts UTF-8 strings from
// protobuf-like binary blobs (tag 0x0a or 0x12, then varint length, then payload).
// Used to present readable text from binary graph/structure blobs.
func ExtractStringsFromBinary(data []byte) []string {
	var out []string
	i := 0
	for i < len(data) {
		if i+2 > len(data) {
			break
		}
		tag := data[i]
		i++
		// Wire type 2 = length-delimited. Tag 1 = 0x0a, tag 2 = 0x12.
		if tag != 0x0a && tag != 0x12 {
			continue
		}
		n, consumed := binary.Uvarint(data[i:])
		if consumed <= 0 {
			break
		}
		i += consumed
		if n == 0 || uint64(len(data)) < uint64(i)+n {
			continue
		}
		chunk := data[i : i+int(n)]
		i += int(n)
		if utf8.Valid(chunk) {
			s := string(chunk)
			if trim := trimSpace(s); trim != "" && len(trim) > 1 {
				out = append(out, trim)
			}
		}
	}
	return out
}

func trimSpace(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
