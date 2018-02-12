1.自定义php.ini
```
ADD ./php.ini /usr/local/etc/php/php.ini
```
2.自定义添加php-fpm.conf
```
ADD ./php-www.conf /usr/local/etc/php-fpm.d/
```
3.启动命令
```
php-fpm -c /usr/local/etc/php/php.ini -y /usr/local/etc/php-fpm.conf
```
