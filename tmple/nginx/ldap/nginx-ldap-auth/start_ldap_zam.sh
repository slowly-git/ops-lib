
#python27安装ldap模块
/home/worker/python27/bin/pip install python-ldap

#登陆页面
su worker -c 'cd /home/worker/nginx/nginx-ldap-auth && nohup /home/worker/python27/bin/python backend-sample-app.py > /dev/null  2>&1 &'
#验证账号密码页面
su worker -c 'cd /home/worker/nginx/nginx-ldap-auth && nohup /home/worker/python27/bin/python nginx-ldap-auth-daemon.py > /dev/null  2>&1 &'