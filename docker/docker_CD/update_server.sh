#!/bin/sh

DOCKER_REG="harbor.tiejin.cn/closer"
DOCKER_IMG="umscloud-server"
DOCKER_IMG_TAG=$1

if [[ ! -n "$1" ]]
then
  DOCKER_IMG_TAG="latest"
else
  DOCKER_IMG_TAG=$1
fi

. $(cd `dirname $0`; pwd)/../basefunction.sh
cp $(cd `dirname $0`; pwd)/../cluster/cluster_production.json /data/${DOCKER_IMG}/conf/cluster.json
source /etc/environment 

docker pull ${DOCKER_REG}/${DOCKER_IMG}:${DOCKER_IMG_TAG}

if [[ $? -eq 0 ]]
then
  docker stop ${DOCKER_IMG}
  sleep 3
  docker rm ${DOCKER_IMG}
else
  echo "docker image not found , exit"
  exit 1
fi

#DOCKER_HOST=`ec2metadata --local-host`
DOCKER_HOST=$(curl http://169.254.169.254/latest/meta-data/local-hostname)
#DOCKER_HOST_IP=`ec2metadata --local-ipv4`
DOCKER_HOST_IP=$(curl http://169.254.169.254/latest/meta-data/local-ipv4)

docker run --name=${DOCKER_IMG} -d -p 8081:8081 -p 8080:8080 -p 8088:8088 -p 8082:8082 -p 8083:8083 -p 8098:8098 -p 9821:9821\
    -e DOCKER_HOST="$DOCKER_HOST" -e DOCKER_HOST_IP="$DOCKER_HOST_IP" \
    -e UMSCLOUD_SERVER_OPTS="-Dserver_name=${DOCKER_IMG} -Denv=production -Dserver_id=$ums_server_id -Dserver_num=$ums_server_num\
    -Dcom.sun.management.jmxremote -Dcom.sun.management.jmxremote.port=8098 -Dcom.sun.management.jmxremote.rmi.port=8098 -Djava.rmi.server.hostname=$DOCKER_HOST_IP -Dcom.sun.management.jmxremote.authenticate=false -Dcom.sun.management.jmxremote.ssl=false \
    -server -Xms20000m -Xmx20000m -Xss512k -Dfile.encoding=utf-8 -XX:MaxMetaspaceSize=100m\
                           -XX:+UseParNewGC -XX:MaxTenuringThreshold=15 -XX:+UseConcMarkSweepGC -XX:SurvivorRatio=13 -XX:CMSInitiatingOccupancyFraction=70 -XX:TargetSurvivorRatio=100\
                           -XX:CMSMaxAbortablePrecleanTime=25000 -XX:+ExplicitGCInvokesConcurrentAndUnloadsClasses -XX:CMSFullGCsBeforeCompaction=5 -XX:+UseCMSCompactAtFullCollection\
                           -XX:+UseCompressedOops -XX:+CMSParallelRemarkEnabled -XX:+CMSScavengeBeforeRemark -XX:+CMSClassUnloadingEnabled -XX:+DisableExplicitGC -XX:ParallelGCThreads=6\
                           -XX:+PrintFlagsFinal -XX:+PrintCommandLineFlags -XX:+PrintGCDateStamps -XX:+PrintTenuringDistribution -XX:+PrintGCDetails -XX:+PrintGCTimeStamps -XX:+PrintGCApplicationStoppedTime\
                           -XX:+PrintGCApplicationConcurrentTime -Xloggc:/logs/server/gc_$time.log \
                           -Dsun.net.inetaddr.ttl=60 \
                           -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=/logs/server/jvm_$time.hprof"\
    -v /data/${DOCKER_IMG}/conf:/data/${DOCKER_IMG}/conf \
    -v /root/.aws:/root/.aws \
    -v /logs/${DOCKER_IMG}/:/logs/server/ ${DOCKER_REG}/${DOCKER_IMG}:${DOCKER_IMG_TAG}
 \

checkserver 127.0.0.1 8080
