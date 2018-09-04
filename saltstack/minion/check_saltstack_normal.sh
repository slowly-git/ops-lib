#!/bin/bash

start_time=$(echo $(stat /dev/zero | grep -E 'Modify|改动' |awk '{print $2,$3,$4}' |xargs -i date -d {} +%s))
now_time=$(date +%s)

dur_time=$(echo $now_time-$start_time|bc -l)

if [[ $dur_time -gt 300 ]]
then

        ###first check hostname
        condition=$(echo $(hostname) |grep 'tiejin.cn' |wc -l)

        if [[ $condition -lt 1 ]]
        then
                echo "xiao yu 1 hostname,will reset hostname"
                cd /home/worker/scripts && ./set_host_name_onstart 

                kill -9 $(ps -ef|grep salt-minion|grep -v grep|awk '{print $2}')

                ### restart salt
                rm -rf /etc/salt/{minion.d,minion_id,pki}
                sleep 2
                /bin/systemctl start salt-minion.service
        fi
        ##CheckSalt-minion is running
        new_salt_running=$(ps -ef |grep salt-minion |grep -v grep |grep "/usr/bin/salt-minion"|awk '{print $8,$9}'|wc -l)

        if [[ $new_salt_running -lt 1 && $condition -eq 1 ]]
        then
                rm -rf /etc/salt/{minion.d,minion_id,pki}
                sleep 2
                /bin/systemctl start salt-minion.service
        fi
fi
