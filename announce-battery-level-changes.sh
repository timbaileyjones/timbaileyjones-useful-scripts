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
  if [ $prev -lt 100 -a $pct -eq 100 ];then
    echo say battery is fully charged at 100%
  fi
  if [ $pct -lt $prev ];then
    echo $(date) $(pmset -g batt | grep InternalBattery )
    if [ $pct -le 10 ];then
      say battery fell from $prev percent to $pct
    fi
  elif [ $pct -gt $prev ];then
    echo $(date) $(pmset -g batt | grep InternalBattery )
    if [ $pct -le 10 ];then
      say battery at $pct %
    fi
  fi
  prev=$pct
  sleep 2
done
