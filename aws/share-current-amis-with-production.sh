#!/bin/bash

aws sts get-caller-identity | grep -q -v 999999999999 
if [ $? -gt 0 ]
then
  echo "Your environment is set up for production.  Set it for development, and re-run."
  exit 1
fi

for KEY in # list of parameter store values
do

  echo -n "adding permission for production to run "
  AMI_ID=`aws ssm --region  $AWS_REGION get-parameter  \
        --name $KEY --output text  --query "Parameter.{Value:Value}"`
  echo -n $AMI_ID 
  aws ec2 --region  $AWS_REGION modify-image-attribute --image-id $AMI_ID \
    --launch-permission "{\"Add\":[{\"UserId\":\"000000000000\"}]}"
  echo .
done
