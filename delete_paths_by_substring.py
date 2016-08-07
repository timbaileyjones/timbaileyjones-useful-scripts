#!/usr/bin/env python 

# vim: number:ts=4:shiftwidth=4:expandtab:nows
    
import re
import os
import sys
import shutil

print "sys.argv", sys.argv

if len(sys.argv) < 3:
    print "Usage: delete_paths_by_substring.py %s %s" % ("filename-list", "regex-to-delete")
    sys.exit(1)

regex = re.compile(sys.argv[2])

filenames = []
with open(sys.argv[1]) as f:
    filenames = f.read().splitlines()

print "filename length", len(filenames)
print "regex", regex

matches = []
nonmatches = []
for f in filenames:
    if regex.search(f):
        print "matched ", f
        matches.append(f)
    else:
        nonmatches.append(f)

print "%d filenames, %d matches" % (len(filenames), len(matches))

del_succ = 0
del_fail = 0
count = 0
for f in matches:
    if (count % 10) == 0:
        print "deleting %d of %d, %.2lf" % (count, len(matches), (100.0*count)/len(matches))
    count = count + 1 
    try:
        os.unlink(f)
        del_succ = del_succ + 1
    except:
        del_fail = del_fail + 1
print

with open("delete-by-regex.log", "a") as f:
    f.write( "%s %d %d\n" % (sys.argv[2], del_succ, del_fail))

with open(".tmp", "w") as f:
    for line in nonmatches:
        f.write(line)
        f.write('\n')

shutil.move(".tmp", sys.argv[1])
