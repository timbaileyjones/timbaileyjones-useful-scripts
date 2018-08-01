#!/bin/bash
aws logs --region $AWS_REGION describe-log-groups | grep logGroupName | grep $1 | cut -f4 -d'"' | while read lg
do
  echo $lg ; aws logs --region $AWS_REGION delete-log-group --log-group-name $lg
done
