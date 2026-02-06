package db

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"google.golang.org/protobuf/encoding/protowire"
)

// We do NOT use the full protobuf library (no .proto schema, no proto.Unmarshal).
// We only:
// 1. Use encoding/protowire to parse the wire format (tag = field number + wire type, then value).
// 2. Heuristically extract UTF-8 from length-delimited fields (ExtractStringsFromBinary).
// Without Cursor's .proto definitions we cannot fully decode; WireDump shows structure
// so we can see what the first bytes represent (e.g. field 1, wire type 2, length 32).

// WireDump parses data as protobuf wire format and returns one line per field:
// "  field N (wire W): ..." with value shown as varint, hex, or UTF-8 preview.
// Uses the official protowire package. If parsing fails, returns a single error line.
func WireDump(data []byte) []string {
	var out []string
	i := 0
	for i < len(data) {
		num, typ, tagLen := protowire.ConsumeTag(data[i:])
		if tagLen < 0 {
			out = append(out, "  (wire parse error at offset "+strconv.Itoa(i)+": "+protowire.ParseError(tagLen).Error()+")")
			break
		}
		rest := data[i+tagLen:]
		valLen := protowire.ConsumeFieldValue(num, typ, rest)
		if valLen < 0 {
			out = append(out, "  field "+strconv.Itoa(int(num))+" (wire "+wireTypeName(typ)+"): (value parse error)")
			break
		}
		payload := rest[:valLen]
		line := "  field " + strconv.Itoa(int(num)) + " (wire " + wireTypeName(typ) + "): "
		switch typ {
		case protowire.VarintType:
			v, _ := protowire.ConsumeVarint(payload)
			line += "varint " + strconv.FormatUint(v, 10)
		case protowire.Fixed32Type:
			v, _ := protowire.ConsumeFixed32(payload)
			line += "fixed32 " + strconv.FormatUint(uint64(v), 10)
		case protowire.Fixed64Type:
			v, _ := protowire.ConsumeFixed64(payload)
			line += "fixed64 " + strconv.FormatUint(v, 10)
		case protowire.BytesType:
			inner, _ := protowire.ConsumeBytes(payload) // length-prefixed; inner is the actual bytes
			if utf8.Valid(inner) && len(inner) > 0 && len(inner) <= 200 {
				s := strings.TrimSpace(string(inner))
				if s != "" {
					line += fmt.Sprintf("bytes %d = %q", len(inner), s)
				} else {
					line += fmt.Sprintf("bytes %d hex %s", len(inner), hex.EncodeToString(inner))
				}
			} else {
				line += fmt.Sprintf("bytes %d hex %s", len(inner), hex.EncodeToString(inner))
				if len(inner) > 32 {
					line += "..."
				}
			}
		case protowire.StartGroupType:
			line += "start_group (nested)"
			// ConsumeFieldValue already advanced over the group; payload is group body
		case protowire.EndGroupType:
			line += "end_group"
		default:
			line += "unknown wire type"
		}
		out = append(out, line)
		i += tagLen + valLen
	}
	return out
}

func wireTypeName(typ protowire.Type) string {
	switch typ {
	case protowire.VarintType:
		return "varint"
	case protowire.Fixed32Type:
		return "fixed32"
	case protowire.Fixed64Type:
		return "fixed64"
	case protowire.BytesType:
		return "bytes"
	case protowire.StartGroupType:
		return "start_group"
	case protowire.EndGroupType:
		return "end_group"
	default:
		return strconv.Itoa(int(typ))
	}
}

// ExtractStringsFromBinary heuristically extracts UTF-8 strings from
// protobuf-like binary blobs (tag 0x0a or 0x12, then varint length, then payload).
// We only look for tags 1 and 2 (length-delimited); protowire would handle all tags.
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
