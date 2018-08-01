#!/bin/bash  
aws ec2 describe-instances \
  --region $AWS_REGION --profile $AWS_PROFILE --output text \
  --filter Name="instance-state-name","Values=running" \
  --query "Reservations[*].Instances[*].{StartTime:LaunchTime,ImageId:ImageId,ZName:Tags[?Key=='Name']|[0].Value,InstanceId:InstanceId,PublicIp:PublicIpAddress,PrivateIp:PrivateIpAddress}" |\
  awk '{ printf "%-20s %-15s %-15s %-15s %-12s %-s\n",$6,$3,$4,$5,$1,$2 }' | sort 

  #--filter Name="key-name",Values="some-keypair" \
