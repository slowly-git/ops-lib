#!/bin/bash

start_time=$(echo $(stat /dev/zero |grep Modify |awk '{print $2,$3,$4}' |xargs -i date -d {} +%s))
now_time=$(date +%s)

dur_time=$(echo $now_time-$start_time|bc -l)

if [[ $dur_time -gt 300 ]]
then

        ###first check hostname
        condition=$(echo $(hostname) |grep 'c360in.com' |wc -l)

        if [[ $condition -lt 1 ]]
        then
                echo "xiao yu 1 hostname,will reset hostname"
                cd /home/worker/scripts && ./set-hostname-onstart && ./create_route53_records
				
		kill -9 $(ps -ef|grep salt-minion|grep -v grep|awk '{print $2}')
				
		### restart old salt
		rm -rf /etc/salt/{minion.d,minion_id,pki}
                sleep 2
                /home/worker/python/bin/salt-minion -c /etc/salt -d
				
		### restart new salt
                rm -rf /usr/local/python27/etc/{minion.d,minion_id}
		rm -rf /etc/salt/pki_new
                sleep 2
                cd /usr/local/python27/bin && ./generate_minion.sh 
        fi
		
	##Check Old-Salt-minion is running
        old_salt_running=$(ps -ef |grep salt-minion |grep -v grep |grep "/home/worker/python/bin/salt-minion"|awk '{print $8,$9}'|wc -l)

        if [[ $old_salt_running -lt 1 && $condition -eq 1 ]]
        then
                rm -rf /etc/salt/{minion.d,minion_id,pki}
                sleep 2
                /home/worker/python/bin/salt-minion -c /etc/salt -d
        fi

        ##Check New-Salt-minion is running
        new_salt_running=$(ps -ef |grep salt-minion |grep -v grep |grep "/usr/local/python27/bin/salt-minion"|awk '{print $8,$9}'|wc -l)

        if [[ $new_salt_running -lt 1 && $condition -eq 1 ]]
        then
                rm -rf /usr/local/python27/etc/{minion.d,minion_id}
		rm -rf /etc/salt/pki_new
                sleep 2
                cd /usr/local/python27/bin && ./generate_minion.sh 
        fi
fi
