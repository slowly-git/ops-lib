#!/bin/bash
### clean expire pods
for i in $(find /data/deployments-*/*/*/*/|awk -F'/' '{print "/"$2"/"$3"/"$4"/"$5"/"$6}'|sort -u)
do
        if [[ ! -f "$i/notice.log" ]]
        then
                echo "$i/notice.log is not exist,will be remove it"
                #touch $i/notice.log
                #echo "$i" |awk -F'/' '{print "/"$2"/"$3"/"$4"/"$5}'
                #echo "$i" |awk -F'/' '{print "/"$2"/"$3"/"$4}'|xargs -i rm -rf {}
        else
                #echo "$i/notice.log is exist,will juge is expire"
                find $i/notice.log -cmin +3 |awk -F'/' '{print "/"$2"/"$3"/"$4"/"$5}' |xargs -i rm -rf {}
        fi
done
