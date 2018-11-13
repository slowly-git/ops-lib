# 自动扩容流程

## 背景
* 服务器根据负载进行伸缩
* 成本节约
* 应对突发流量

## 现状
目前因为业务初期，线上为固定的机器提供服务

## 自动扩容设计
利用aws的autoscaling+cloudwathc服务，对有伸缩可能的服务进行服务改造，实现自动扩容，自动提供服务，且无需人工介入
主要用到的组件包括:cloudwatch,autoscaling,saltstack
* cloudwatch: aws自带的监控服务，主要使用其cpu的监控
* autoscaling: aws的服务，提供服务器自动伸缩的功能
* saltstack: Linux集群管理系统，保证自动扩容后服务的正常启动
## 流程
成员 | 角色
--------- | -------------
cloudWatch |  aws自有监控
autoscaling | aws自动伸缩服务
salt master | Linux集群管理中心
ec2 | 伸缩的业务机
ELB | aws负载均衡器
```sequence
    cloudWatch->autoscaling:触发扩容
    autoscaling->ec2:ec2扩容
    ec2->ec2:执行开机脚本
    ec2->salt master:注册
    salt master-->ec2:分发部署脚本
    ec2->ec2:部署业务
    ec2->ELB:注册服务器到elb
    ELB-->ec2:健康检查后转发流量
    
```


