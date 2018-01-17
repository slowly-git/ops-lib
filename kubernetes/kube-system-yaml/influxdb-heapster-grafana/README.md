# grafana
1.指定node

给node打标签

```
kubectl label nodes 10.100.2.112 task=monitor
```

2.指定挂载/data目录，子目录为grafana-storage

influxdb相似

3.grafana自带的对k8s influxdb以及cluster.json 和 pods.json，因为在设置了grafana密码后，不能自动生成，
需要手动配置influxDB以及cluster和pods的dashboard,相应配置参见:

官方github:https://github.com/kubernetes/heapster/tree/master/grafana/dashboards

或者 dashboard目录下的json文件
