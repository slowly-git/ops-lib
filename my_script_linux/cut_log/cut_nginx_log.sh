#!/usr/bin/env bash

LOGFILE=`basename ${0}`
exec 1>> `dirname ${0}`/${LOGFILE}.log
exec 2>> `dirname ${0}`/${LOGFILE}.log
 
LOGS_DIR="$HOME/data/nginx/logs"
DAY_DIR=$LOGS_DIR/$(date -d "yesterday" +"%Y")/$(date -d "yesterday" +"%m")/$(date -d "yesterday" +"%d")
DAY_GZIP_DIR=$LOGS_DIR/$(date -d "-2 day" +"%Y")/$(date -d "-2 day" +"%m")/$(date -d "-2 day" +"%d")
DATA_DIR=$LOGS_DIR/`date +"%Y"`
timeago=3
 
if [ ! -d $DAY_DIR ]
then
  mkdir -p $DAY_DIR
fi
mv $LOGS_DIR/*access.log* $DAY_DIR/
mv $LOGS_DIR/*error.log* $DAY_DIR/

sudo /home/worker/nginx/sbin/nginx -s reopen

find $DATA_DIR -type f -mtime +${timeago} |xargs rm -v
find $LOGS_DIR -empty -type d |xargs rm -vrf