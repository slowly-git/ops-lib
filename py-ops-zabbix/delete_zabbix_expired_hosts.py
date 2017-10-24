#!/usr/bin/env python
# coding=utf-8
# feature: delete zabbix hosts which is not health
# 删除自动注册机器中已经下线的机器

import json, urllib2
import logging
import os


ZABBIX_API_URL = "https://en-us-zabbix.360in.com/zabbix/api_jsonrpc.php"
ZABBIX_HEADER = {"Content-Type": "application/json"}
ZABBIX_USER = "apicall"
ZABBIX_PASS = "apicall@999"


def _log(filename):
    log_level = logging.INFO
    logger = logging.getLogger("delete-zabbix-hosts")
    logger.setLevel(log_level)
    fh = logging.FileHandler(os.path.join('/tmp', filename))
    fh.setLevel(log_level)
    ch = logging.StreamHandler()
    ch.setLevel(log_level)
    formatter = logging.Formatter("%(asctime)s [%(levelname)s]-%(name)s: %(message)s")
    fh.setFormatter(formatter)
    ch.setFormatter(formatter)
    logger.addHandler(fh)
    logger.addHandler(ch)
    return logger


def url_request(post_data, log):
    """
    公共请求抽离
    :param log:
    :param post_data:
    :return:
    """
    request = urllib2.Request(ZABBIX_API_URL, post_data)
    for key in ZABBIX_HEADER:
        request.add_header(key, ZABBIX_HEADER[key])
    try:
        result = urllib2.urlopen(request, timeout=300)
    except urllib2.URLError as e:
        log.error(str(e))
        return False
    else:
        response = json.loads(result.read())
        result.close()
        return response


def get_auth_code(log):
    """
    获取到zabbix的登录授权码
    :return:
    """
    auth_para = {
        "jsonrpc": "2.0",
        "method": "user.login",
        "params": {
            "user": ZABBIX_USER,
            "password": ZABBIX_PASS
        },
        "id": 1
    }
    response = url_request(json.dumps(auth_para), log)
    try:
        assert isinstance(response, dict)
    except AssertionError:
        log.error("auth code return type should be dict.")
        return False

    if response.has_key('error'):
        log.error(response['error'])
        return False

    if response['result']:
        log.info("get auth code: {}".format(response['result']))
        return response['result'] # this is auth code


def get_unhealthy_hosts(auth_code, log):
    """
    获取非正常状态的主机
    :param auth_code:
    :param log:
    :return:
    """
    r_data = {
        "jsonrpc": "2.0",
        "method": "host.get",
        "params": {
            "output": ["hostid", "name", "available", "maintenance_status", "status", "error"],
            "limit": 20,
            "filter": {
                "available": ["2"],
           #     "status": ["0"]
                # status: 0,监控中； 1,未监控
                # available：0,可用未知； 1,可用；2,不可用
                # "host": "en-ap-mq-app-20160721143713-de99d2d6acd1110568895ba"
            },
        },
        "auth": auth_code,
        "id": 1
    }
    response = url_request(json.dumps(r_data), log)
    return response['result']


def delete_unhealthy_hosts(unhealthy_hosts, auth_code, log):
    """
    删除错误的主机
    :param log:
    :param unhealthy_hosts:
    :param auth_code:
    :return:
    """
    host_ids = []
    host_names = []
    for host in unhealthy_hosts:
        if "cannot connect" in  host['error']:
            host_names.append(host['name'])
            host_id = host['hostid']
            host_ids.append(host_id)

    log.info("delete unhealthy hosts: {}".format(host_names))
    log.info("delete unhealthy hosts count is: {}".format(len(host_ids)))
    print host_ids
    if len(host_ids) > 0:
        delete_data = {
        "jsonrpc": "2.0",
        "method": "host.delete",
        "params": host_ids,
        "auth": auth_code,
        "id": 1
    }
        response = url_request(json.dumps(delete_data), log)
        if response.has_key('error'):
            log.error("delete failed: {}".format(response['error']))
        else:
            log.info("delete zabbix hosts success.")
    else:
        log.info("there is no matched host to delete, ignore")


def main():
    log = _log("delete-zabbix-hosts.log")
    auth_code = get_auth_code(log)
    unhealthy_hosts = get_unhealthy_hosts(auth_code, log)
#    print unhealthy_hosts
    delete_unhealthy_hosts(unhealthy_hosts, auth_code, log)


if __name__ == '__main__':
    main()