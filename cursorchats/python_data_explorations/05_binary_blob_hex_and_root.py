"""
05 - Binary blobs: root blob (latestRootBlobId) is binary; short blobs show hex.
We see 0x0a often (protobuf tag for field 1, length-delimited).
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
cur.execute("SELECT value FROM meta WHERE key = '0'")
row = cur.fetchone()
meta = json.loads(bytes.fromhex(row[0]).decode("utf-8"))
root_id = meta["latestRootBlobId"]
print("Root blob id (latestRootBlobId):", root_id)

cur.execute("SELECT id, length(data), data FROM blobs WHERE id = ?", (root_id,))
row = cur.fetchone()
if row:
    bid, length, d = row
    print("Root blob len:", length)
    print("First 100 bytes (hex):", d[:100].hex())
    print("Valid UTF-8:", d.decode("utf-8", errors="replace")[:80] if d else "n/a")

# One short binary blob (find in Python to avoid SQL on blob bytes)
cur.execute("SELECT id, length(data), data FROM blobs LIMIT 100")
for row in cur.fetchall():
    bid, length, d = row[0], row[1], row[2]
    if length and length < 100 and d and not d.startswith(b'{"'):
        print("\nShort binary blob id:", bid[:24], "len:", length)
        print("Hex:", d.hex())
        break
conn.close()
