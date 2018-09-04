### saltstack minion会使用到的脚本
```
### watch salt minion is running
*/1 * * * * /bin/bash -x /home/worker/scripts/check_saltstack_normal.sh > /home/worker/scripts/check_salt_minion_alive.log 2>&1
```
