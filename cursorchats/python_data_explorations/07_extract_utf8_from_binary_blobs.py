"""
07 - Heuristic: extract UTF-8 strings from binary blobs (protobuf-like layout).
Many blobs use 0x0a (tag 1, length-delimited) then one-byte length then payload.
This recovers readable text from binary blobs for display.
"""
import os
import sqlite3

def extract_strings(data):
    i = 0
    out = []
    while i < len(data):
        if data[i] == 0x0A and i + 1 < len(data):  # tag 1, length-delimited
            n = data[i + 1]
            i += 2
            if n > 0 and i + n <= len(data):
                chunk = data[i : i + n]
                i += n
                try:
                    s = chunk.decode("utf-8")
                    if s.strip():
                        out.append(s.strip())
                except Exception:
                    pass
            continue
        if data[i] == 0x12 and i + 1 < len(data):  # tag 2, length-delimited
            n = data[i + 1]
            i += 2
            if n > 0 and i + n <= len(data):
                chunk = data[i : i + n]
                i += n
                try:
                    s = chunk.decode("utf-8")
                    if len(s.strip()) > 2:
                        out.append(s.strip())
                except Exception:
                    pass
            continue
        i += 1
    return out


base = os.path.expanduser("~/.cursor/chats")
db_path = None
for root, _dirs, files in os.walk(base):
    for f in files:
        if f == "store.db":
            db_path = os.path.join(root, f)
            break
    if db_path:
        break
if not db_path:
    print("No store.db found"); exit(1)

conn = sqlite3.connect(db_path)
cur = conn.cursor()
# Get blobs that are not JSON (no leading "{")
cur.execute("SELECT id, data FROM blobs")
rows = cur.fetchall()
conn.close()

count = 0
for row in rows:
    bid, data = row[0], row[1]
    if not data or data.startswith(b'{"'):
        continue
    if len(data) > 500:
        continue
    strings = extract_strings(data)
    if not strings:
        continue
    print("BLOB id=%s len=%d" % (bid[:24], len(data)))
    for s in strings[:6]:
        print("   ", repr(s[:180]))
    print()
    count += 1
    if count >= 5:
        break
