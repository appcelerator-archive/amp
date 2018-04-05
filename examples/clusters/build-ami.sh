#!/bin/bash
DIR="$(dirname $0)"
cd "$DIR"
DIR="$(pwd -P)"
image="appcelerator/ansible:2.5"
REGIONS="us-east-1 us-east-2 us-west-2 eu-west-1 ap-southeast-2"
OWNER=${OWNER:-654814900965}
IMAGE_NAME=${IMAGE_NAME:-ubuntu-xenial-docker}

# check for credentials in env variables, else check credential files
if [[ -n "$AWS_ACCESS_KEY_ID" && -n "$AWS_SECRET_ACCESS_KEY" ]]; then
  _cred_opts="-e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY"
elif [[ -f $HOME/.aws/credentials || -f $HOME/.aws/config ]]; then
  _cred_opts="-v $HOME/.aws:/root/.aws:ro"
else
  echo "Please configure your aws credentials first"
  exit 1
fi

# aws cli binary or image
which aws &>/dev/null
if [[ $? -ne 0 ]]; then
  echo "No aws CLI available, will use a Docker image"
  aws="docker run --rm $_cred_opts cgswong/aws:aws"

else
  aws="aws"
fi
$aws --version

# check that the credentials are from the expected account
if [[ -n "$AWS_ACCESS_KEY_ID" && -n "$AWS_SECRET_ACCESS_KEY" ]]; then
  account=$($aws iam get-user | grep Arn | cut -d: -f6)
  if [[ $? -ne 0 ]]; then
    exit 1
  fi
  if [[ "$account" != "$OWNER" ]]; then
    echo "credentials are from account $account, expected $OWNER, abort"
    exit 1
  fi
fi

# build the AMI
docker run --rm -v "$DIR:/data" $_cred_opts "$image" ansible-playbook /data/build-ami.yml || exit 1

echo "update the cloudformation template mapping for key 'Default' with the values below:"
for region in $REGIONS; do
  printf "${region}: "
  $aws --region "$region" ec2 describe-images --owners $OWNER --filters Name=name,Values="${IMAGE_NAME}*" --query 'Images[*].{ID:ImageId,Date:CreationDate}' --output text  | sort | tail -1 | awk '{print $NF}'
done
