import getpass
import random, base64
from hashlib import sha1

user = getpass.getuser()
passwd = getpass.getpass(prompt="请输入密码")
# if svc_login(user, passwd):  # You must write svc_login()
#     print('ok')
# else:
#     print('Boo!')
