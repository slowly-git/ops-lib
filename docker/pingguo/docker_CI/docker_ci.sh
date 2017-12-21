#!/bin/env bash

COMMAND_NAME=`basename $0`
if [ $# -lt 2 ]
then
  echo "Usage: ${COMMAND_NAME} projectname dockerfile"
  exit 1
fi

IMAGE_REG="en-us-public-ops-harbor-1.360in.com:4333/pinguo/"
IMAGE_TAG=`date  +%Y-%m%d-%H%M-%S`
IMAGE_PROJECT=$1
IMAGE_NAME=${IMAGE_REG}${IMAGE_PROJECT}-${IMAGE_TAG}

IMAGE_DOCKERFILE=$2

# build docker images
$(cd ${IMAGE_DOCKERFILE} && /usr/bin/docker build --no-cache -t ${IMAGE_NAME} . > /dev/null)

# CD in kubernetes
/usr/local/bin/kubectl set image deployments/ads-ad-k8s ads-ad-k8s: --namespace=ads

# return to jenkins
echo "${IMAGE_NAME} build done!"