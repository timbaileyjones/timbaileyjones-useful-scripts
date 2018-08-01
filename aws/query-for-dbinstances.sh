#!/bin/bash  
aws rds describe-db-instances \
  --region $AWS_REGION --profile $AWS_PROFILE --output text \
  --query "DBInstances[*].{DBInstanceIdentifier:DBInstanceIdentifier,Endpoint:Endpoint.Address}" | pr -te10 | column -c 1
