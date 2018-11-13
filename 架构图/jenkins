# 发布流程6-25
## 背景
* 线上业务需要规范自动发布
* 发布过后需要有回滚机制
* 发布机制需要兼容autlscaling
* 权限管理

## 现状
现在使用的teamcity+ssh的方式发布代码到固定的几台服务器上，如果出现大量的服务器缩容和扩容，无法满足自动化的发布需求，且没有相关的回滚机制
 
## 建议采用的发布流程架构
主要用到的组件包括：jenkins,saltstack,harbor
* jenkins:主要负责获取业务镜像并发布到线上服务器
* saltstack:主要负责对Linux服务器集群进行配置管理，满足服务的自动伸缩
* harbor：线上镜像的存储仓库，配合AWS S3服务使用

## 流程
#### 权限管理
成员 | 角色
--------- | -------------
user |  项目发布负责人
Jenkins | 发布中心 
ldap | 用户验证
auth policy | jenkins项目授权策略
```sequence
    user->jenkins:用户名密码登陆
    jenkins->ldap:用户认证
    ldap-->jenkins:用户认证通过
    jenkins->auth policy:根据用户获取相应授权
    auth policy-->user:用户获取指定项目权限
```
####  jenkins发布流程
成员 | 角色
--------- | -------------
user |  项目发布负责人
Jenkins | 发布中心 
salt master | 配置管理服务端
salt minion | 配置管理客户端
harbor | docker镜像仓库
```sequence
    user->jenkins:登陆
    jenkins-->user:登陆成功
    user->job:点击构建
    job->harbor:获取版本号
    harbor-->jenkins:返回版本号
    user->job:点击发布
    job->job: pull new img
    job->job: tag latest img
    job->harbor: push img
    harbor-->job: push down
    job->salt master: 下发发布命令
    salt master->salt minion:执行awselb命令
    salt minion->salt minion:
    Note left of salt minion: 更新容器
    salt minion-->salt master:回传发布结果
    salt master-->jenkins: console 展示
    jenkins-->user:发布完成
```
#### awselb命令执行流程（自动化）
成员 | 角色
--------- | -------------
ec2 | 服务器
salt minion |  配置管理客户端
elb | aws负载均衡器
awselb| aws elb 管理工具(运维开发)
bashshell | 业务更新脚本
```sequence
    salt minion->awselb:执行更新命令
    awselb->elb: 从ELB中摘除ec2
    elb-->awselb: 摘除成功
    awselb->bashshell: 调用更新脚本更新服务
    bashshell-->awselb: 服务更新完成
    awselb->elb: 注册ec2到elb并循环检查监控状况
    elb-->awselb: 状态检查通过（不通过则阻塞，告警）
    awselb-->salt minion:返回执行结果
```
#### bashshell脚本执行流程（自动化）
成员 | 角色
--------- | -------------
bashshell | 业务更新脚本
tag | 容器镜像版本号
docker server | 业务容器
harbor | docker镜像仓库
```sequence
bashshell->tag:获取版本号
tag-->bashshell:如果tag为空，则默认为latest
bashshell->harbor:下载docker img
harbor-->bashshell:done
bashshell->docker server:删除重启
docker server-->bashshell:done

```




