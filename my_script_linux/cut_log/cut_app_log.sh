#!/bin/bash

BASE_DIR=/home/worker/data/www/runtime
TIME=`date +'%Y%m%d%H' -d '-1 hours'`
OLDDAY=`date +'%Y%m%d' -d '-5 days'`
WHITELIST='cloud_fileLogProcess.log'
WHITESYS='user'
TIMEAGO=1


LOGFILE=`basename ${0}`
exec 1>> /home/worker/scripts/${LOGFILE}.log
exec 2>> /home/worker/scripts/${LOGFILE}.log

function hitwhite() {
        local f 
        local ret=0
        for f in $WHITELIST
        do
                if [[ "X$f" = "X$1" ]]; then
                        ret=1
                        break
                fi
        done
        echo $ret
}

function hitwhitesys() {
        local f 
        local ret=0
        for f in $WHITESYS
        do
                if [[ "X$f" = "X$1" ]]; then
                        ret=1
                        break
                fi
        done
        echo $ret
}

date
if [[ -d $BASE_DIR ]]; then
for dir in `ls $BASE_DIR`
do 
        #echo $dir
        ignoresys=`hitwhitesys $dir`
        if [[ $ignoresys -eq 1 ]]; then
                continue;
        fi
        for file in `find $BASE_DIR/$dir/ -maxdepth 8 -name '*.log'`
        do
                filename=`basename $file`
                ignore=`hitwhite $filename`
                if [[ -f $file && $ignore -eq 0 ]]; then
                        mv -v $file $file.$TIME
                fi
        done
        for file in `find $BASE_DIR/$dir/ -maxdepth 8 -name "*.log.${OLDDAY}*"`
        do 
                if [[ -f $file ]]; then
                        rm -v $file
                fi
        done
        # delete notice 1 ago file
        for file in `find $BASE_DIR/$dir/ -type f -ctime +${TIMEAGO} -maxdepth 8 -name "notice.log.*"`
        do
                if [[ -f $file ]]; then
                        rm -v $file
                fi
        done
        # delete cloud_log 3 ago file
        if [[ "$dir" == "cloud" ]]
        then
                find $BASE_DIR/$dir/cloud_log -type f -ctime +3 |xargs rm -v
                find $BASE_DIR/$dir/cloud_log -empty -type d |xargs rm -vrf
        fi
done
fi