# coding:utf-8
'''
zabbix.py：公用zabbix模块
_zabbix_url：zabbix地址
_zabbix_usr：zabbix用户名
_zabbix_passwd：zabbix密码
'''
import json, urllib.request
from urllib.request import URLError
import sys, argparse

_zabbix_url = 
_zabbix_usr = 
_zabbix_passwd = 


class zabbix_api:
    def __init__(self):
        self.url = _zabbix_url
        self.header = {"Content-Type": "application/json"}

    def user_login(self):
        '''
        获取用于认证zabbix的authID
        :return: self.authID
        '''
        data = json.dumps({
            "jsonrpc": "2.0",
            "method": "user.login",
            "params": {
                "user": _zabbix_usr,  # 修改用户名
                "password": _zabbix_passwd  # 修改密码
            },
            "id": 0
        })
        request = urllib.request.Request(self.url, bytes(data, 'utf8'))
        for key in self.header:
            request.add_header(key, self.header[key])
        try:
            result = urllib.request.urlopen(request)
        except URLError as e:
            print("用户认证失败，请检查用户名和密码", e.reason)
        else:
            response = json.loads(result.read())
            result.close()
            self.authID = response['result']
            return self.authID

    def host_get(self, hostName=""):
        '''
        获取监控主机
        :param hostName: 要查询的主机名称，默认为空，查询所有信息
        :return host_info: 返回主机信息字典
        host_info["host_count"] = 监控数量
        host_info["host_name"] = 监控主机名称
        host_info["host_search"] = 查询主机状态
        '''
        data = json.dumps({
            "jsonrpc": "2.0",
            "method": "host.get",
            "params": {
                "output": "extend",
                "filter": {
                    "host": "",
                    "name": hostName
                }
            },
            "auth": self.user_login(),
            "id": 1
        })
        request = urllib.request.Request(self.url, data.encode('utf8'))
        for key in self.header:
            request.add_header(key, self.header[key])
        try:
            result = urllib.request.urlopen(request)
        except URLError as e:
            # 判断报错e对象是否有'reason'属性
            if hasattr(e, 'reason'):
                print('''
                登录失败！
                失败原因:{_reason}
                '''.format(_reason=e.reason))
            elif hasattr(e, 'code'):
                print('''
                服务器无响应！
                错误码:{_code}
                '''.format(_code=e.code))
        else:
            response = json.loads(result.read())
            result.close()
            host_info = dict()
            host_info["host_count"] = len(response['result'])
            host_info["host_name"] = []
            host_info["host_search"] = {}
            for host in response['result']:
                if len(hostName) == 0:
                    host_info["host_name"].append(host['name'])
                else:
                    host_info["host_search"] = host
            return host_info


if __name__ == "__main__":
    zabbix = zabbix_api()
    parser = argparse.ArgumentParser(description='zabbix  api ', usage='%(prog)s [options]')
    parser.add_argument('-H', '--host', nargs='?', dest='listhost', default='host', help='查询主机')
    parser.add_argument('-v', '--version', action='version', version='%(prog)s 1.0')
    if len(sys.argv) == 1:
        print(parser.print_help())
    else:
        args = parser.parse_args()
        if args.listhost != 'host':
            if args.listhost:
                host = zabbix.host_get(hostName=args.listhost)['host_search']
                print(host)
            else:
                host = zabbix.host_get()['host_name']
                print(host)
