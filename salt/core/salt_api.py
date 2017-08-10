import json
import urllib.request
from urllib.request import URLError
import ssl


class SaltApi(object):
    def __init__(self, **kwargs):
        self.salt_url = kwargs['url']  # salt 地址
        self.salt_user = kwargs['user']
        self.salt_password = kwargs['password']
        self.header = {"Content-Type": "application/json"}
        self.context = ssl.SSLContext(ssl.PROTOCOL_TLS)
        self.context.verify_mode = ssl.CERT_REQUIRED
        self.context.check_hostname = False
        self.context.load_verify_locations('/etc/pki/tls/certs/localhost.crt')

    def uer_login(self):
        """
        - 获取token
        """
        data = json.dumps({
            'eauth': 'pam',
            'username': self.salt_user,
            'password': self.salt_password
        }).encode('utf-8')
        request = urllib.request.Request(self.salt_url + '/login', data)
        for key in self.header:
            request.add_header(key, self.header[key])
        try:
            result = urllib.request.urlopen(request, context=self.context)
        except URLError as e:
            print("salt用户认证失败，请检查用户名和密码", e.reason)
        else:
            response = json.loads(result.read())
            result.close()
            auth_token = response['return'][0]['token']
            return auth_token

    def list_all_key(self):
        headers = self.header
        headers.update({'X-Auth-Token': self.uer_login()})
        data = json.dumps({
            'client': 'wheel',
            'fun': 'key.list_all'
        }).encode('utf-8')
        request = urllib.request.Request(self.salt_url, data)
        for key in headers:
            request.add_header(key, headers[key])
        try:
            result = urllib.request.urlopen(request, context=self.context)
        except URLError as e:
            print("salt用户认证失败，请检查用户名和密码", e.reason)
        else:
            response = json.loads(result.read().decode('utf8'))
            result.close()
            minions = response['return'][0]['data']['return']['minions']
            return minions

    def get_minion(self):
        pass
