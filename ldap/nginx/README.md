# nginx ldap

##### install nginx

```
./configure --prefix=/data/nginx \
--with-http_realip_module \
--with-http_ssl_module \
--with-http_sub_module \
--with-http_auth_request_module \
--with-http_stub_status_module
```
##### get nginx-ldap-auth
```
cd /home/worker/nginx
git clone https://github.com/nginxinc/nginx-ldap-auth

```
##### install python-ldap
```
/home/worker/python27/bin/pip install python-ldap
```

#####   启动登陆前端和认证程序
```
# start nginx
su worker -c 'cd /home/worker/nginx/nginx-ldap-auth && nohup /home/worker/python27/bin/python backend-sample-app.py > /dev/null  2>&1 &'
su worker -c 'cd /home/worker/nginx/nginx-ldap-auth && nohup /home/worker/python27/bin/python nginx-ldap-auth-daemon.py > /dev/null  2>&1 &'
```
##### 配置nginx.conf
```
#缓存可以减少ldap验证频率，不然每个页面都需要ldap验证一次
proxy_cache_path cache/ keys_zone=auth_cache:10m;
```
##### 配置vhost
```
server
{
    listen 8080;
    server_name logs.tiejin.cn;
    location / {
        auth_request /auth-proxy;

        #nginx接收到nginx-ldap-auth-daemon.py返回的401和403都会重新跳转到登录页面
        error_page 401 403 =200 /login;

        proxy_pass http://kibana_logcenter;
    }
   access_log logs/kibana_logcenter_access.log;

    #登录页面，由backend-sample-app.py提供
    location /login {
        proxy_pass http://127.0.0.1:9000/login;
        proxy_set_header X-Target $request_uri;
    }

    location = /auth-proxy {
        internal;
        proxy_pass http://127.0.0.1:8888;     #nginx-ldap-auth-daemon.py运行端口

        #缓存设置
        proxy_cache auth_cache;
        proxy_cache_key "$http_authorization$cookie_nginxauth";
        proxy_cache_valid 200 403 10m;
        
        proxy_pass_request_body off;
        proxy_set_header Content-Length "";

        #这些配置都会通过http头部传递给nginx-ldap-auth-daemon.py脚本
        proxy_set_header X-Ldap-URL      "ldap://cn-bj-public-ops-freeipa01.tiejin.cn:389";
        proxy_set_header X-Ldap-BaseDN   "cn=users,cn=accounts,dc=tiejin,dc=cn";
        proxy_set_header X-Ldap-BindDN   "uid=nginx,cn=users,cn=accounts,dc=tiejin,dc=cn";
        proxy_set_header X-Ldap-BindPass "xxxxxxxxxx";

        proxy_set_header X-Ldap-Template "(uid=%(username)s)";

        proxy_set_header X-CookieName "nginxauth";
        proxy_set_header Cookie nginxauth=$cookie_nginxauth;
    }
}
```


