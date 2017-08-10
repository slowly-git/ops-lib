#!/usr/bin/python
# -*-coding: utf-8 -*-

import os, socket
import datetime

HOSTNAME = socket.gethostname()
HOME = os.environ['HOME']

APP_LOG_FILE = HOME + '/data/var/report/monitor/bdp.log'
ZABBIX_CONF = HOME + "/zabbix/etc/zabbix_agentd.conf"
DATA_FILE = HOME + "/data/zabbix/bdp_stat.data"
LOG_FILE = HOME + "/data/zabbix/bdp_stat.log"

nowtime = datetime.datetime.now()
deltime = datetime.timedelta(minutes=-1)
log_minu = nowtime + deltime
GET_TIME = log_minu.strftime('%Y/%m/%d %H:%M')

#GET_TIME = os.popen('''date -d '-1 min' +'%Y/%m/%d  %H:%M:%S'|awk '{print $2}'|awk -F ':' '{position=$1":"$2;print position}' ''').read().strip()

os.environ['APP_LOG_FILE'] = str(APP_LOG_FILE)
os.environ['ZABBIX_CONF'] = str(ZABBIX_CONF)
os.environ['DATA_FILE'] = str(DATA_FILE)
os.environ['LOG_FILE'] = str(LOG_FILE)
os.environ['GET_TIME'] = str(GET_TIME)

log = os.popen('''/bin/grep -R "$GET_TIME" $APP_LOG_FILE''').read()

if log:
    file_obj = open(DATA_FILE, 'w+')
    for x in log.split('\n')[:-1]:
        file_obj.write(HOSTNAME + '\t' + x.split()[2] + '\t' + x.split()[3] + '\n')
    file_obj.close()
    os.popen("~/bin/zabbix_sender -vv -c $ZABBIX_CONF -i $DATA_FILE >>$LOG_FILE 2>&1")
else:
    file_obj = open(DATA_FILE, 'a')
    file_obj.write("{0},no data found! \n".format(GET_TIME))
    file_obj.close()
