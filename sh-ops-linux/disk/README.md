#### 磁盘性能测试
##### 1.安装lib包
```
$ yum install libaio.x86_64 libaio-devel.x86_64 -y
``` 

##### 2.安装fio工具
```
$ wget http://brick.kernel.dk/snaps/fio-2.1.10.tar.gz 
$ tar -zxf fio-2.1.10.tar.gz
$ cd fio-2.1.10
$ ./configure
$ make && make install
```

##### 3.运行脚本
