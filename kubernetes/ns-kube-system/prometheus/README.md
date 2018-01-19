# prometheus
1.指定node

给node打标签

```
kubectl label nodes 10.100.2.112 task=monitor
```

2.指定挂载/data目录，子目录为prometheus-storage

3.给相应目录权限，通过设置
```
  securityContext:
    runAsUser: 0
```

