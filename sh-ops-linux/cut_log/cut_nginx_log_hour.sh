#!/bin/bash

BASE_DIR="${HOME}/data/nginx/logs"
TIME=`date +'%Y%m%d%H' -d '-1 hours'`
NGINX_STATUS="$(ps -ef|grep 'nginx:'|wc -l)"

for file in `find $BASE_DIR/ -maxdepth 1 -name '*.log'`
do
  mv -v $file $file.$TIME
  sudo /home/worker/nginx/sbin/nginx -s reopen
done
