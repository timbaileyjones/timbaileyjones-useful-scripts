#!/bin/bash  
aws rds describe-db-clusters \
  --region $AWS_REGION --profile $AWS_PROFILE --output text \
  --query "DBClusters[*].{DatabaseName:DatabaseName,Endpoint:Endpoint}" | column -c 1
