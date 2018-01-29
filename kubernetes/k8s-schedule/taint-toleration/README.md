Taint（污点）和 Toleration（容忍）可以作用于 node 和 pod 上，其目的是优化 pod 在集群间的调度，
这跟节点亲和性类似，只不过它们作用的方式相反，具有 taint 的 node 和 pod 是互斥关系，而具有节
点亲和性关系的 node 和 pod 是相吸的。另外还有可以给 node 节点设置 label，通过给 pod 设置 
nodeSelector 将 pod 调度到具有匹配标签的节点上。

Taint 和 toleration 相互配合，可以用来避免 pod 被分配到不合适的节点上。每个节点上都可以应用一个或多个 taint ，
这表示对于那些不能容忍这些 taint 的 pod，是不会被该节点接受的。如果将 toleration 应用于 pod 上，
则表示这些 pod 可以（但不要求）被调度到具有相应 taint 的节点上

```
operator可以定义为：

Equal 表示key是否等于value，默认
Exists 表示key是否存在，此时无需定义value


effect可以定义为：

NoSchedule 表示不允许调度，已调度的不影响
PreferNoSchedule 表示尽量不调度
NoExecute 表示不允许调度，已调度的在tolerationSeconds（定义在Tolerations上）后删除
```


#为 node1 设置 taint:

创建taint:
```
kubectl taint nodes node1 key1=value1:NoSchedule
kubectl taint nodes node1 key1=value1:NoExecute
kubectl taint nodes node1 key2=value2:NoSchedule

```

删除taint:
```
kubectl taint nodes node1 key1:NoSchedule-
kubectl taint nodes node1 key1:NoExecute-
kubectl taint nodes node1 key2:NoSchedule-

```

查看taint:
```
kubectl describe nodes node1 
```

#为 pod 设置 toleration

```
tolerations:
- key: "key1"
  operator: "Equal"
  value: "value1"
  effect: "NoSchedule"
- key: "key1"
  operator: "Equal"
  value: "value1"
  effect: "NoExecute"
- key: "node.alpha.kubernetes.io/unreachable"
  operator: "Exists"
  effect: "NoExecute"
  tolerationSeconds: 6000
```
value 的值可以为 NoSchedule、PreferNoSchedule 或 NoExecute。
tolerationSeconds 是当 pod 需要被驱逐时，可以继续在 node 上运行的时间。


