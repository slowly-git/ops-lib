### saltstack服务器会使用到的脚本
```
###delete_dead_salt_minion
*/2 * * * * /root/scripts/delete_salt_dead_machine.sh >>/dev/shm/deleted_minion
###restart salt master & api
0 3 * * *  /root/scripts/restart_salt.sh > /dev/null 2>&1

```
### saltstack minion会使用到的脚本
```
### watch salt minion is running
*/1 * * * * /bin/bash -x /home/worker/scripts/check_saltstack_normal.sh > /home/worker/scripts/check_salt_minion_alive.log 2>&1
```


