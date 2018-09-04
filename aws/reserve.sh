#!/bin/bash

aws ec2 describe-reserved-instances |grep ^RESERVEDINSTANCES|awk '$(NF-1)~"active"{print $6,$8}' >/dev/shm/value
cat /dev/shm/value  |awk '{print $2}'|sort -u >/dev/shm/key

for i in $(cat /dev/shm/key)
do
        echo -ne $i,$(grep -w "$i" /dev/shm/value|awk 'BEGIN{sum=0}{sum+=$1}END{print sum}') "\n"
done|sort -t',' -k2nr >/dev/shm/reserve_final


echo "#################预留数量比正在使用的数量多几个##############################"

aws ec2 describe-instances |grep  -E -w  "^INSTANCES|^STATE" |grep -B1 running |grep -w ^INSTANCES|awk '{for (i=1; i<=NF; i++) if ($i ~ /(nano|micro|small|medium|large|xlarge|2xlarge|4xlarge|8xlarge)/) print $i}'|uniq -c >/dev/shm/value_online

#aws ec2 describe-instances |grep  -E -w  "^INSTANCES|^STATE" |grep -B1 running |grep -w ^INSTANCES|grep -E "2018|jumper02|closer"|awk '{print $8}'|uniq -c >>/dev/shm/value_online

cat /dev/shm/value_online |awk '{print $2}'|sort -u >/dev/shm/key_online

for i in $(cat /dev/shm/key_online)
do
        echo -ne $i,$(grep -w "$i" /dev/shm/value_online|awk 'BEGIN{sum=0}{sum+=$1}END{print sum}') "\n"
done|sort -t',' -k2nr >/dev/shm/online_final


for i in $(cat /dev/shm/key_online)
do
        reserve_num=$(grep -w $i /dev/shm/reserve_final|wc -l)
        online_num=$(grep -w $i /dev/shm/online_final|wc -l)
        if [[ $reserve_num -ge 0 && $online_num -ge 1 ]]
        then
                echo -ne $i && echo -en "\t" && echo $(grep -w $i /dev/shm/reserve_final |awk -F',' '{print $2}') - $(grep -w $i /dev/shm/online_final|awk -F',' '{print $2}') |bc -l
        fi
done

for i in $(cat /dev/shm/key)
do
        reserve_num=$(grep -w $i /dev/shm/reserve_final|wc -l)
        online_num=$(grep -w $i /dev/shm/online_final|wc -l)
        if [[ $reserve_num -ge 1 && $online_num -eq 0 ]]
        then
                echo -ne $i && echo -en "\t" && echo $(grep -w $i /dev/shm/reserve_final |awk -F',' '{print $2}') - 0 |bc -l
        fi
done

cd /dev/shm/ && rm -rf *value* *final* *key*
