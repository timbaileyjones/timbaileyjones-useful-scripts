"""
03 - Blobs table: mix of JSON (chat messages) and binary.
Classify by whether data starts with '{"' and is valid JSON with role/content.
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
cur.execute("SELECT id, length(data), data FROM blobs")
rows = cur.fetchall()
conn.close()

json_count = binary_count = 0
sample_json = None
for row in rows:
    _id, length, d = row[0], row[1], row[2]
    if d and d.startswith(b'{"'):
        try:
            j = json.loads(d.decode("utf-8"))
            if j.get("role") and j.get("content"):
                json_count += 1
                if sample_json is None:
                    sample_json = (d, j)
        except Exception:
            binary_count += 1
    else:
        binary_count += 1

print("Blobs: JSON (role+content):", json_count, "  binary/other:", binary_count)
if sample_json:
    d, j = sample_json
    print("\nSample JSON blob (first 500 chars of content):")
    print(j.get("content", "")[:500])
