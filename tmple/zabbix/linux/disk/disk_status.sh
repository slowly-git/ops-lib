#!/bin/env bash

source `dirname $0`/zabbixConfig.sh

export PATH=$PATH:$HOME/bin
COMMAND_NAME=`basename $0`
if [ $# -lt 1 ]
then
  echo "Usage: $COMMAND_NAME device"
  exit 1
fi
DEVICE=$1

HOSTNAME=`hostname`
DATA_FILE="$HOME/data/zabbix/${COMMAND_NAME}_$HOSTNAME_$DEVICE.data"
LOG_FILE="$HOME/data/zabbix/${COMMAND_NAME}_$HOSTNAME_$DEVICE.log"

STATUS=`iostat -kx 1 2| grep "$DEVICE " | tail -n 1`

echo "$HOSTNAME disk.$DEVICE.iostat.wrqm `echo "$STATUS" | awk '{print $2}'`" >$DATA_FILE
echo "$HOSTNAME disk.$DEVICE.iostat.wrqm `echo "$STATUS" | awk '{print $3}'`" >>$DATA_FILE
echo "$HOSTNAME disk.$DEVICE.iostat.r `echo "$STATUS" | awk '{print $4}'`" >>$DATA_FILE
echo "$HOSTNAME disk.$DEVICE.iostat.w `echo "$STATUS" | awk '{print $5}'`" >>$DATA_FILE
echo "$HOSTNAME disk.$DEVICE.iostat.rkB `echo "$STATUS" | awk '{print $6}'`" >>$DATA_FILE
echo "$HOSTNAME disk.$DEVICE.iostat.wkB `echo "$STATUS" | awk '{print $7}'`" >>$DATA_FILE
echo "$HOSTNAME disk.$DEVICE.iostat.avgrq-sz `echo "$STATUS" | awk '{print $8}'`" >>$DATA_FILE
echo "$HOSTNAME disk.$DEVICE.iostat.avgqu-sz `echo "$STATUS" | awk '{print $9}'`" >>$DATA_FILE
echo "$HOSTNAME disk.$DEVICE.iostat.await `echo "$STATUS" | awk '{print $10}'`" >>$DATA_FILE
echo "$HOSTNAME disk.$DEVICE.iostat.svctm `echo "$STATUS" | awk '{print $11}'`" >>$DATA_FILE
echo "$HOSTNAME disk.$DEVICE.iostat.util `echo "$STATUS" | awk '{print $12}'`" >>$DATA_FILE

~/bin/zabbix_sender -vv -c $zabbix_agentd_config -i $DATA_FILE >>$LOG_FILE 2>&1

echo 0

