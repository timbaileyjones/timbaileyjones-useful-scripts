"""
04 - Chronological order: meta['createdAt'] (ms) vs file mtime.
We use createdAt to sort chats; mtime as fallback when meta is missing.
"""
import json
import os
import sqlite3

base = os.path.expanduser("~/.cursor/chats")
paths = []
for root, _dirs, files in os.walk(base):
    for f in files:
        if f == "store.db":
            paths.append(os.path.join(root, f))

print("Path (relative)                    mtime       createdAt (from meta)")
print("-" * 80)
for path in sorted(paths)[:15]:
    rel = path.replace(base, "")
    mtime = int(os.path.getmtime(path))
    created = None
    name = ""
    try:
        conn = sqlite3.connect(path)
        cur = conn.cursor()
        cur.execute("SELECT value FROM meta WHERE key = '0'")
        row = cur.fetchone()
        conn.close()
        if row:
            raw = bytes.fromhex(row[0]).decode("utf-8")
            data = json.loads(raw)
            created = data.get("createdAt")
            name = (data.get("name") or "")[:25]
    except Exception:
        pass
    print(f"{rel[:40]:40} {mtime}  {created}  {name}")
