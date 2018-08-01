aws cloudformation --region $AWS_REGION --output text list-exports  | awk '{ print $3"\t"$4 }' | pr -te40 | column -c 1 | sort 
