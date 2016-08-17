#!/bin/bash
# vim:ai:ts=4:shiftwidth=4:

REPOS="repo1 repo2 repo3"
if [ -f setenv.sh ]
then
    . setenv.sh 
fi
mkdir -p workspaces

set -e
for repo in $REPOS
do
	#title svnadmin $repo
	if [ ! -d $repo ]
	then
		svnadmin create $repo
		echo '#!/bin/sh' > $repo/hooks/pre-revprop-change
		chmod 755 $repo/hooks/pre-revprop-change
		svnsync init file://$PWD/$repo $SVN_URL/$repo 
	fi 
	#title svnsync $repo
	svnsync sync file://$PWD/$repo 
done

for repo in $REPOS
do
	if [ ! -d workspaces/$repo ]
	then
		cd workspaces
		#title svn checkout $repo
		svn checkout $SVN_URL/$repo 
		cd ../
	else
		cd workspaces/$repo
		#title svn update $repo
		svn update
		cd ../../
	fi
done
exit 0
