#elasticsearch

```
image: gcr.io/google-containers/elasticsearch:v5.6.4
dockerfile: https://github.com/kubernetes/kubernetes/tree/master/cluster/addons/fluentd-elasticsearch/es-image
```
1. 需要设置系统内核参数
```
vm.max_map_count=26214
```
如果其他镜像，则要在yml文件通过initcontainer设置
```
initContainers:
- image: alpine:3.6
  command: ["/sbin/sysctl", "-w", "vm.max_map_count=262144"]
  name: elasticsearch-logging-init
  securityContext:
    privileged: true
```

2.此处使用两个elasticsearch集群



3.集群测试
```
[root@cn-office-ops-zam01 tmp]# curl "http://172.30.37.5:9200/_cat/health"         
1516861706 06:28:26 docker-cluster yellow 1 1 16 16 0 0 16 0 - 50.0%
```
查看数据情况
```
[root@cn-office-ops-zam01 tmp]# curl -XGET "http://172.30.37.5:9200/_cat/indices?v"
health status index                       uuid                   pri rep docs.count docs.deleted store.size pri.store.size
yellow open   .monitoring-es-6-2018.01.25 GFTqesnrSG6MpwPyA5c5Lw   1   1       9939          160      4.4mb          4.4mb
yellow open   logstash-2018.01.25         8pBzs7BHTH2ifZOkx2J3gA   5   1       9176            0      6.1mb          6.1mb
yellow open   logstash-2018.01.23         rhi1S9cYSX2WVamRVIgJVg   5   1        178            0    151.2kb        151.2kb
yellow open   logstash-2018.01.24         oeZIM3T8Qo6ezPWsrAR2mQ   5   1      10315            0      3.5mb          3.5mb
```
