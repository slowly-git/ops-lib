#!/bin/env bash
# $1 项目名称，会用来以此命名镜像
# $2 dockerfile 的全路径
# $3....n k8s deployment项目名称
#
LOCK_FILE=/tmp/docker_ci_lock_file

if [[ ! -f ${LOCK_FILE} ]]
then  
        touch ${LOCK_FILE}
else
        echo "${COMMAND_NAME} is working now"
        exit 1
fi  

COMMAND_NAME=`basename $0`
if [ $# -lt 3 ] 
then
  echo "Usage: ${COMMAND_NAME} projectname dockerfile k8sdeployment"
  rm -rf ${LOCK_FILE}
  exit 1
fi

IMAGE_REG="en-us-public-ops-harbor-1.360in.com:4333/pinguo/"
IMAGE_TAG=`date  +%Y-%m%d-%H%M-%S`
IMAGE_PROJECT=$1
IMAGE_NAME=${IMAGE_REG}${IMAGE_PROJECT}-${IMAGE_TAG}

IMAGE_DOCKERFILE=$2

#DEPLOY_NAME=$3
DEPLOY_LIST=$@

# build docker images
$(cd ${IMAGE_DOCKERFILE} && /usr/bin/docker build --no-cache -t ${IMAGE_NAME} . > /dev/null)


# push image
/usr/bin/docker push ${IMAGE_NAME}


# update deployment

if [[ $? -eq 0 ]]
then
  for DEPLOY_NAME in ${DEPLOY_LIST[@]:3}
  do
    su - root -c "kubectl set image deployment/${DEPLOY_NAME} ${DEPLOY_NAME}=${IMAGE_NAME} -n ads"
  done
else
        echo "cant update deployment"
fi

rm -rf ${LOCK_FILE}
