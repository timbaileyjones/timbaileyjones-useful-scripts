#!/bin/bash
#

function batt() {
    pmset -g batt | grep InternalBattery | awk '{ print $3 }' | cut -f1 -d%
}

prev=$(batt)    # grab current battery level as previous, so that it doesn't announce
                # battery status every single time you start this script.

while :
do
  pct=$(batt)
  if [ $pct -lt $prev ];then
    echo $(date) $(pmset -g batt | grep InternalBattery )
    if [ $pct -le 30 ];then
      say battery fell from $prev percent to $pct
    fi
  elif [ $pct -gt $prev ];then
    echo $(date) $(pmset -g batt | grep InternalBattery )
    if [ $pct -le 30 ];then
      say battery at $pct %
    fi
  fi
  prev=$pct
  sleep 10
done
