# metrics

1,默认的，插件 api server 会引用 kube-system 中的configmap:extension-apiserver-authentication

但有些插件需要我们自己手动指定: metrics

```
command:
- --requestheader-client-ca-file=/var/run/kubernetes/client-ca-file

volumeMounts:
- name: client-ca
  mountPath: /var/run/kubernetes

volumes:
  - name: client-ca
    configMap:
      name: extension-apiserver-authentication
```

2,测试api是否正常
```
 kubectl get --raw "/apis/metrics.k8s.io/v1beta1/nodes"|jq .
```

3,metrics-server替代了hepster,后续需要结合prometheus+custom-metrics-apiserver自定义监控

4. k8s deploy doc:https://github.com/kubernetes/kubernetes/tree/master/cluster/addons/metrics-server
