# 简介

flunted收集k8s系统日志到es

filebeat通过side-car模式收集app日志到kafka，再由logstash消费到es

kibana展示日志


1. 建立namespace log
```
kubectl create namespace log

```

2.启动kafka，然后进入容器查看状态
```
[root@cn-office-ops-zam01 ns-log]# kubectl exec -it kafka-with-zook-1-5c6657cbbb-cstwp /bin/bash -c kafka1 -n log   

```

如果打开了kafka JMX监控，进入容器后需要先去除环境变量JMX_PORT
```
bash-4.3# unset JMX_PORT
```

查看topic是否创建
```
bash-4.3# cd /opt/kafka/bin
bash-4.3#./kafka-topics.sh --zookeeper 127.0.0.1:2181 --list
k8s

```
查看topic状态
```
bash-4.3# ./kafka-topics.sh --zookeeper 127.0.0.1:2181 --topic "k8s" --describe                  
Topic:k8s       PartitionCount:3        ReplicationFactor:1     Configs:
        Topic: k8s      Partition: 0    Leader: 1       Replicas: 1     Isr: 1
        Topic: k8s      Partition: 1    Leader: 1       Replicas: 1     Isr: 1
        Topic: k8s      Partition: 2    Leader: 1       Replicas: 1     Isr: 1
```



