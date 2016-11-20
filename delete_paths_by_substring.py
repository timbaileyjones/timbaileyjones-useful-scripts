#!/usr/bin/env python 

# vim: number:ts=4:shiftwidth=4:expandtab:nows
    
import re
import os
import sys
import shutil

print "sys.argv", sys.argv

if len(sys.argv) < 3:
    print "Usage: delete_paths_by_substring.py %s %s" % ("filename-list", "file-of-regexes-to-delete")
    sys.exit(1)



regexes=[]
patterns2delete=[]
with open(sys.argv[2]) as f:
    patterns2delete = f.read().splitlines()
    for pattern in patterns2delete:
        if len(pattern) > 8:
            r = re.compile(pattern)
            regexes.append(r)

filenames = []
with open(sys.argv[1]) as f:
    filenames = f.read().splitlines()

print "filename length", len(filenames)

matches = []
nonmatches = []

for f in filenames:
    found = False
    for i, r in enumerate(regexes):
        if r.search(f):
            #print "matched regex %s: %s" % (patterns2delete[i], f)
            matches.append(f)
            found = True
            break
    if not found:
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

#
#  log the deletion stats for this run
#
with open("delete-by-regex.log", "a") as f:
    f.write( "%s %d %d\n" % (patterns2delete.__str__(), del_succ, del_fail))

with open(".tmp", "w") as f:
    for line in nonmatches:
        f.write(line)
        f.write('\n')
shutil.move(".tmp", sys.argv[1])

#
# rewrite patterns-to-delete file, skipping the patterns we just deleted
#
old_patterns2delete=[]
with open(sys.argv[2]) as f:
    old_patterns2delete = f.read().splitlines()
with open(".tmp", "w") as f:
    for line in old_patterns2delete:
        if line not in patterns2delete:
            f.write(line)
            f.write('\n')

shutil.move(".tmp", sys.argv[2])
