#!/bin/env bash

export PATH=$PATH:$HOME/bin

source `dirname $0`/zabbixConfig.sh

COMMAND_NAME=`basename $0`
HOSTNAME=`hostname`
DATA_FILE="$HOME/data/zabbix/${COMMAND_NAME}_$HOSTNAME.data"
LOG_FILE="$HOME/data/zabbix/${COMMAND_NAME}_$HOSTNAME.log"
date=
month=`date +%Y%m`
day=`date +%d`
tm=`date +%Y%m%d%H`
date_php=$(date +%d-%b-%Y" "%H:%M -d'-1 minute')
date_time=$(date +%Y/%m/%d" "%H:%M -d'-1 minute')
short_application_name=$(echo $1|awk -F'.' '{print $1}')
fls=$(echo $1)

logdata()
{
	##########For php_fpm make error_log############
	if [[ -f /home/worker/data/php/log/php_errors.log ]]
	then
		php_error_log_num=$(grep "$date_php" /home/worker/data/php/log/php_errors.log |wc -l)
	else
		php_error_log_num=0
	fi

	##########for application nginx error_log##############
	if [[ -f $HOME/data/nginx/logs/"$fls".error.log ]]
	then
		application_nginx_error_log_num=$(grep "$date_tim" $HOME/data/nginx/logs/"$fils".error.log |wc -l)
	else
		application_nginx_error_log_num=0
	fi

	############for application runtime error log ################
	#if [[ -f $HOME/data/www/runtime/$short_application_name/error.log ]]
	#then
	#	application_runtime_error_log=$(grep "$date_time" $HOME/data/www/runtime/$short_application_name/error.log |wc -l)
	#else
	#	application_runtime_error_log=0
	#fi
	###########for application runtime error log ################
	if [[ -f $HOME/data/www/runtime/$short_application_name/application.log ]]
	then
		application_runtime_error_log=$(grep "$date_time" $HOME/data/www/runtime/$short_application_name/application.log |wc -l)
	else
		application_runtime_error_log=0
	fi

	###########for php_fpm slow log##################
	if [[ -f $HOME/data/php/log/www.log.slow ]]
	then
		application_php_slow_log_num=$(grep "$date_php" $HOME/data/php/log/www.log.slow |wc -l)
	else
		application_php_slow_log_num=0
	fi

}                        




zabbixdata()
{
  
   if [[ "$php_error_log_num" -gt 0 ]]
   then
        echo "$HOSTNAME php_error_log_num $php_error_log_num" >$DATA_FILE
   else
        echo "$HOSTNAME php_error_log_num 0" >$DATA_FILE
   fi

   if [[ "$application_nginx_error_log_num" -gt 0 ]]
   then
        echo "$HOSTNAME application_nginx_error_log_num $application_nginx_error_log_num" >>$DATA_FILE
   else
        echo "$HOSTNAME application_nginx_error_log_num 0" >>$DATA_FILE
   fi

   if [[ "$application_runtime_error_log" -gt 0 ]]
   then
        echo "$HOSTNAME application_runtime_error_log $application_runtime_error_log" >>$DATA_FILE
   else
        echo "$HOSTNAME application_runtime_error_log 0" >>$DATA_FILE
   fi

   if [[ "$application_php_slow_log_num" -gt 0 ]]
   then
        echo "$HOSTNAME application_php_slow_log_num $application_php_slow_log_num" >>$DATA_FILE
   else
        echo "$HOSTNAME application_php_slow_log_num 0" >>$DATA_FILE
   fi

   zabbix_sender -vv -c $zabbix_agentd_config -i $DATA_FILE >>$LOG_FILE 2>&1


}


rm -rf $LOG_FILE
logdata
zabbixdata
