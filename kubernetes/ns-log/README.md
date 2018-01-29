
1.镜像来源

es,flunted为k8s插件中的版本:
```
地址：https://github.com/kubernetes/kubernetes/tree/master/cluster/addons/fluentd-elasticsearch

es:      gcr.io/google-containers/elasticsearch:v5.6.4
flunted: gcr.io/google-containers/fluentd-elasticsearch:v2.0.4

```
kibana,logstash,filebeat为elastic官方版本：

```
地址：https://www.docker.elastic.co/#

docker.elastic.co/kibana/kibana:5.6.4
docker.elastic.co/logstash/logstash:6.1.2
docker.elastic.co/beats/filebeat:6.1.2    //kafka output插件还不支持kafka1.0版本

```
