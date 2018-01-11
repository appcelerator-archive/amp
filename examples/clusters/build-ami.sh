#!/bin/bash
DIR="$(dirname $0)"
cd "$DIR"
DIR="$(pwd -P)"
image="appcelerator/ansible:2.4"
REGIONS="us-east-1 us-east-2 us-west-2 eu-west-1 ap-southeast-2"
OWNER=654814900965
IMAGE_NAME=ubuntu-xenial-docker

if [[ ! -f $HOME/.aws/credentials && ! -f $HOME/.aws/config ]]; then
  echo "Please configure your aws credentials first"
  exit 1
fi
docker run --rm -v "$DIR:/data" -v "$HOME/.aws:/root/.aws:ro" "$image" ansible-playbook /data/build-ami.yml || exit 1

which aws &>/dev/null
if [[ $? -ne 0 ]]; then
  echo "aws cli is not available, check the output above to get the ami IDs and update the cloudformation template" >&2
  exit 0
fi
echo "update the cloudformation template mapping for key 'Default' with the values below:"
for region in $REGIONS; do
  printf "${region}: "
  aws --region "$region" ec2 describe-images --owners $OWNER --filters Name=name,Values="${IMAGE_NAME}*" --query 'Images[*].{ID:ImageId,Date:CreationDate}' --output text  | sort | tail -1 | awk '{print $NF}'
done
