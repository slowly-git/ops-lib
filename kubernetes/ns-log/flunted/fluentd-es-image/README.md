# 因为要给fluntd-es-image添加kafka插件，所以需要自己重做镜像

1. 给相应脚本执行权限
```
chmod a+x clean-apt
chmod a+x clean-install
chmod a+x clean-install

```

2. 制作镜像

```
docker build -t harbor.past123.com/past123/fluentd-elasticsearch:1.0 .
```
