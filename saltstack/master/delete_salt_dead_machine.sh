#!/bin/bash
timeout 30 salt --hide-timeout --no-color --out-file /root/good '*'  test.ping
cat /root/good |grep ':'|awk -F':' '{print $1}' >/root/live_minion
salt-key -L |grep -Ev "Keys" >/root/all_minion

grep -vFf /root/live_minion /root/all_minion |xargs -i salt-key -d {} -y
