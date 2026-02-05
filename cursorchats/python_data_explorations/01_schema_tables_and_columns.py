"""
01 - Discover the schema: list tables and their columns.
Run this first to see that store.db has two tables: meta and blobs.
"""
import os
import sqlite3

# Use first store.db found under ~/.cursor/chats
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
    print("No store.db found under", base)
    exit(1)

print("Database:", db_path)
conn = sqlite3.connect(db_path)
cur = conn.cursor()

cur.execute("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
tables = cur.fetchall()
print("TABLES:", [t[0] for t in tables])

for (name,) in tables:
    cur.execute(f"PRAGMA table_info({name})")
    cols = cur.fetchall()
    print(f"\n--- {name} ---")
    for c in cols:
        print(" ", c)

conn.close()
