#!/bin/bash
if [ $# -lt 3 ]
then
    echo This script is intended to be invoked by copy-ami-ids-from-development.sh.
    exit 1
fi

export KEY=$1
export AWS_ACCESS_KEY_ID=$2
export AWS_SECRET_ACCESS_KEY=$3
export AWS_SESSION_TOKEN=$4
unset AWS_PROFILE
aws ssm --region  $AWS_REGION get-parameter  \
        --name $KEY --output text  --query "Parameter.{Value:Value}"
