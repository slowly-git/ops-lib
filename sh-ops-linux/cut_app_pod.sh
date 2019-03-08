#!/bin/bash

BASE_DIR='/data/deployments-*'
TIME=`date +'%Y%m%d%H' -d '-1 hours'`
TIMEAGO=240

#roate log
for file in $(find $BASE_DIR -maxdepth 8 -name '*.log')
do
        if [[ -f $file ]]; then
                mv -v $file $file.$TIME && touch $file
        fi
done

# delete notice TIMEAGO  file
for file in $(find $BASE_DIR -maxdepth 8 -type f -cmin +${TIMEAGO} -name 'notice.log.*')
do
        if [[ -f $file ]]; then
                rm -v $file
        fi
done
# delete application TIMEAGO  file
for file in $(find $BASE_DIR -maxdepth 8 -type f -cmin +${TIMEAGO} -name 'application.log.*')
do
        if [[ -f $file ]]; then
                rm -v $file
        fi
done
