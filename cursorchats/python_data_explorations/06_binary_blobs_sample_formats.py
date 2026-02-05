"""
06 - Sample binary and non-role/content JSON blobs.
Shows hex prefix (0a = protobuf length-delimited field 1) and JSON with other keys.
"""
import json
import os
import sqlite3

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
cur.execute("SELECT id, data FROM blobs")
rows = cur.fetchall()

print("=== JSON blobs that are NOT role+content (other keys) ===")
for row in rows:
    bid, data = row[0], row[1]
    if not data or not data.startswith(b'{"'):
        continue
    try:
        j = json.loads(data.decode("utf-8"))
        if j.get("role") and j.get("content"):
            continue
        print("JSON_OTHER id=%s len=%d keys=%s" % (bid[:20], len(data), list(j.keys())))
        print("  sample:", json.dumps(j)[:300])
        print()
    except Exception:
        pass

print("=== Binary blobs (first 32 bytes hex) ===")
count = 0
for row in rows:
    bid, data = row[0], row[1]
    if not data or data.startswith(b'{"'):
        continue
    print("BINARY id=%s len=%d hex=%s" % (bid[:20], len(data), data[:32].hex()))
    count += 1
    if count >= 10:
        break
conn.close()
