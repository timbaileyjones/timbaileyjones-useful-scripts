#!/bin/bash
export TMP=/gdrive-lo1
set -e

for vm in `vboxmanage list vms | cut -f2 -d\{ | cut -f1 -d\} | egrep -v 'ca2c6867-d39b-46db-bd0c-be4200739cf3|0bc5b421-20dd-4d5a-bdb2-c30c8008bae6|774ce518-0b19-4b70-8f4d-e04202b10497'`
do
	freespace1=`df -k . | grep $PWD  | awk '{ print $4 }'`

	name=`vboxmanage list vms | grep $vm | cut -f2 -d\" `
	echo -n checking to see if $name has been exported ... 
	if [ ! -f $TMP/$vm.ova ]
	then
		echo not yet. exporting OVA ...
		time vboxmanage export $vm --output $TMP/$vm.ova
	else
		echo yes.  Skipping exporting.
	fi
	echo deleting $name ... 
	time vboxmanage unregistervm  $vm --delete
	echo importing $name ... 
	time vboxmanage import  $TMP/$vm.ova && rm -f $TMP/$vm.ova && 

	freespace2=`df -k . | grep $PWD  | awk '{ print $4 }'`
	echo started with $freespace1 free, now we have $freespace2 free.  Freed up `echo $freespace2 - $freespace1 | bc `.
	echo ========================================================================================================
done


