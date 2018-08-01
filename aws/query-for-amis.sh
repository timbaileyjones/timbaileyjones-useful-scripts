#!/bin/bash
aws --region $AWS_REGION --profile $AWS_PROFILE --output text ec2 describe-images  \
  --query "Images[*].{ID:ImageId,Name:Name,CreateDate:CreationDate,Public:Public}"  \
  --filters "Name=name,Values=ac-*" | grep -v 'True$' |\
   sort | cut -f1-3 | tail -4 | column -c 1
