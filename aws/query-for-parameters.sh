#!/bin/bash  
parameters=`aws ssm describe-parameters \
  --region $AWS_REGION --profile $AWS_PROFILE --output text  \
  --query "Parameters[*].Name" `

# sort the parameters in order (space->LF, sort, then LF->space)
parameters=`echo $parameters | tr " " "\n" | sort | tr "\n" " " `



for parameter in $parameters
do 
  aws ssm get-parameter --region $AWS_REGION --profile $AWS_PROFILE --name $parameter   \
  --output text  --query "Parameter.{Name:Name,Value:Value}" #| column -c 1

done
