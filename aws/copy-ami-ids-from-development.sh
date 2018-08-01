#!/bin/bash
#
# new role name in dev: 
#
# rolename: read-only-access-to-dev-paramstore
# arn: arn:aws:iam::000000000000:role/read-only-access-to-dev-paramstore
#
#new policy in production:
#   assume-role-for-dev-paramstore

aws sts get-caller-identity | grep -q 000000000000 
if [ $? -gt 0 ]
then
  echo "Your environment isn't set up for production.  Fix that, and re-run"
  exit 1
fi

ROLE_ARN=arn:aws:iam::345369269725:role/read-only-access-to-dev-paramstore
SESSION_NAME=dev

# Switch to the role with a new session
ROLE_DATA=$(aws sts assume-role --role-arn "$ROLE_ARN" --role-session-name "$SESSION_NAME" --output text --query 'Credentials.[AccessKeyId,SecretAccessKey,SessionToken]')

#
# parsing out the development keys, to be passed to the helper script for development
#
export DEV_AWS_ACCESS_KEY_ID=$(echo $ROLE_DATA | awk '{print $1}')
export DEV_AWS_SECRET_ACCESS_KEY=$(echo $ROLE_DATA | awk '{print $2}')
export DEV_AWS_SESSION_TOKEN=$(echo $ROLE_DATA | awk '{print $3}')

echo -n retrieving the AMI IDs from development...
CUSTOMER_AMI_api=`./get-ami-from-development.sh CUSTOMER_AMI_api $DEV_AWS_ACCESS_KEY_ID $DEV_AWS_SECRET_ACCESS_KEY $DEV_AWS_SESSION_TOKEN`
CUSTOMER_AMI_app=`./get-ami-from-development.sh CUSTOMER_AMI_app $DEV_AWS_ACCESS_KEY_ID $DEV_AWS_SECRET_ACCESS_KEY $DEV_AWS_SESSION_TOKEN`
CUSTOMER_AMI_awr=`./get-ami-from-development.sh CUSTOMER_AMI_awr $DEV_AWS_ACCESS_KEY_ID $DEV_AWS_SECRET_ACCESS_KEY $DEV_AWS_SESSION_TOKEN`
CUSTOMER_AMI_SFTP_DROPSERVER=`./get-ami-from-development.sh CUSTOMER_AMI_SFTP_DROPSERVER $DEV_AWS_ACCESS_KEY_ID $DEV_AWS_SECRET_ACCESS_KEY $DEV_AWS_SESSION_TOKEN`

echo 
echo "  CUSTOMER_AMI_api=$CUSTOMER_AMI_api"
echo "  CUSTOMER_AMI_app=$CUSTOMER_AMI_app"
echo "  CUSTOMER_AMI_awr=$CUSTOMER_AMI_awr"
echo "  CUSTOMER_AMI_SFTP_DROPSERVER=$CUSTOMER_AMI_SFTP_DROPSERVER"

echo writing the AMI IDs to production...
echo "  CUSTOMER_AMI_api:" `aws ssm --region $AWS_REGION put-parameter --name CUSTOMER_AMI_api             --type String --overwrite --value $CUSTOMER_AMI_api `
echo "  CUSTOMER_AMI_app:" `aws ssm --region $AWS_REGION put-parameter --name CUSTOMER_AMI_app             --type String --overwrite --value $CUSTOMER_AMI_app `
echo "  CUSTOMER_AMI_awr:" `aws ssm --region $AWS_REGION put-parameter --name CUSTOMER_AMI_awr             --type String --overwrite --value $CUSTOMER_AMI_awr `
echo "  CUSTOMER_AMI_SFTP_DROPSERVER:" `aws ssm --region $AWS_REGION put-parameter --name CUSTOMER_AMI_SFTP_DROPSERVER --type String --overwrite --value $CUSTOMER_AMI_SFTP_DROPSERVER `


exit 1
