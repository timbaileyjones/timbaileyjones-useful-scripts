#!/bin/bash  
aws elbv2 describe-load-balancers \
  --region $AWS_REGION --profile $AWS_PROFILE --output text \
  --query "LoadBalancers[*].{LoadBalancerName:LoadBalancerName,DNSName:DNSName,Scheme:Scheme}" |\
   pr -te45 | column -c 1
