参考：https://hub.docker.com/_/php/

一、在项目中使用php-cli

1、构建dockerfile
```
FROM php:7.0-cli
COPY . /usr/src/myapp
WORKDIR /usr/src/myapp
CMD [ "php", "./your-script.php" ]
```
2、启动docker镜像
```
$ docker build -t my-php-app .
$ docker run -it --rm --name my-running-app my-php-app
```
二、直接用php-cli执行一个脚本
```
$ docker run -it --rm --name my-running-script -v "$PWD":/usr/src/myapp -w /usr/src/myapp php:7.0-cli php your-script.php
```
三、安装PHP扩展
扩展安装工具
```
docker-php-ext-configure
docker-php-ext-install
docker-php-ext-enable
```
1、核心扩展安装用例
```
FROM php:7.0-fpm
RUN apt-get update && apt-get install -y \
        libfreetype6-dev \
        libjpeg62-turbo-dev \
        libmcrypt-dev \
        libpng-dev \
    && docker-php-ext-install -j$(nproc) iconv mcrypt \
    && docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/ \
    && docker-php-ext-install -j$(nproc) gd
```
2、PECL 扩展安装用例
```
FROM php:7.1-fpm
RUN pecl install redis-3.1.0 \
    && pecl install xdebug-2.5.0 \
    && docker-php-ext-enable redis xdebug
FROM php:5.6-fpm
RUN apt-get update && apt-get install -y libmemcached-dev zlib1g-dev \
    && pecl install memcached-2.2.0 \
    && docker-php-ext-enable memcached
```
3、其他扩展
```
FROM php:5.6-apache
RUN curl -fsSL 'https://xcache.lighttpd.net/pub/Releases/3.2.0/xcache-3.2.0.tar.gz' -o xcache.tar.gz \
    && mkdir -p xcache \
    && tar -xf xcache.tar.gz -C xcache --strip-components=1 \
    && rm xcache.tar.gz \
    && ( \
        cd xcache \
        && phpize \
        && ./configure --enable-xcache \
        && make -j$(nproc) \
        && make install \
    ) \
    && rm -r xcache \
    && docker-php-ext-enable xcache
```