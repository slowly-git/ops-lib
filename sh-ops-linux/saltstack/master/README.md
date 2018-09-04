###saltstack服务器会使用到的脚本
```
###delete_dead_salt_minion
*/2 * * * * /root/scripts/delete_salt_dead_machine.sh >>/dev/shm/deleted_minion
###restart salt master & api
0 3 * * *  /root/scripts/restart_salt.sh > /dev/null 2>&1

```

