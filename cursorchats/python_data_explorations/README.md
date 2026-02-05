# Python data explorations

Numbered scripts showing how our understanding of the Cursor chat `store.db` schema evolved. Run in order (01 → 07).

**Prerequisite:** `~/.cursor/chats` must exist and contain at least one `store.db` (i.e. you have used Cursor chat).

| # | Script | What it shows |
|---|--------|----------------|
| 01 | `01_schema_tables_and_columns.py` | Tables: `meta` and `blobs`; column names and types |
| 02 | `02_meta_hex_json_decode.py` | Meta key `"0"` is hex-encoded JSON (agentId, createdAt, name, etc.) |
| 03 | `03_blobs_json_vs_binary.py` | Blobs: some JSON with `role`/`content`, rest binary |
| 04 | `04_chronological_createdAt_vs_mtime.py` | Using `createdAt` from meta for ordering; mtime as fallback |
| 05 | `05_binary_blob_hex_and_root.py` | Root blob (latestRootBlobId) is binary; 0x0a suggests protobuf |
| 06 | `06_binary_blobs_sample_formats.py` | Sample binary hex and JSON blobs with other keys |
| 07 | `07_extract_utf8_from_binary_blobs.py` | Heuristic: extract UTF-8 strings from protobuf-like binary blobs |

Run all (from repo root):

```bash
cd cursorchats/python_data_explorations
for f in 0*.py; do echo "=== $f ==="; python3 "$f"; echo; done
```
