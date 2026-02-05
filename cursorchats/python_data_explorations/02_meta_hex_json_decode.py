"""
02 - Meta table: key '0' holds hex-encoded JSON.
Decode it to see agentId, latestRootBlobId, name, createdAt, etc.
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
cur.execute("SELECT key, value FROM meta WHERE key = '0'")
row = cur.fetchone()
conn.close()

if not row:
    print("No meta row with key '0'"); exit(1)

key, value = row
# Value is hex-encoded JSON
decoded = bytes.fromhex(value).decode("utf-8")
data = json.loads(decoded)
print("Meta key 0 (decoded):")
print(json.dumps(data, indent=2))
