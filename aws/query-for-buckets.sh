#!/bin/bash  
aws s3 ls \
  --region $AWS_REGION --profile $AWS_PROFILE --output text \
#  --query "Reservations[*].Instances[*].{StartTime:LaunchTime,ImageId:ImageId,ZName:Tags[?Key=='Name']|[0].Value,InstanceId:InstanceId,PublicIp:PublicIpAddress,PrivateIp:PrivateIpAddress}" |\
#  awk '{ printf "%-20s %-15s %-15s %-15s %-12s %-s\n",$6,$3,$4,$5,$1,$2 }' | sort 
