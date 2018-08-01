#!/bin/bash  
if [ 0 = 1 ]
then
  aws rds describe-db-cluster-snapshots --region $AWS_REGION --profile $AWS_PROFILE --output json 
fi

for region in $AWS_REGION us-west-2
do 
  echo snapshots in $region 
  aws rds describe-db-cluster-snapshots \
    --region $region --profile $AWS_PROFILE --output text \
    --query "DBClusterSnapshots[*].{ASnapshotCreateTime:SnapshotCreateTime,Status:Status,DBClusterSnapshotIdentifier:DBClusterSnapshotIdentifier,ARN:DBClusterSnapshotArn}" |\
  sort | pr -te32 | column -c 1
done
