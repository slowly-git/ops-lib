#!/usr/bin/python
# coding: utf-8
# author: chenlin@camera360.com
# zabbix 通过微信企业号来推送告警， 主要解决
# （1）微信access token缓存问题，编码反复请求生成token影响发送效率
# （2）zabbix告警串行执行，重复发送相同告警的问题
import hashlib
import urllib2
import json
import sys
import memcache
import logging
import os
reload(sys)
sys.setdefaultencoding('utf8')


def log(filename):
    log_level = logging.INFO
    logger = logging.getLogger("send-weixin")
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


class MC(object):
    def __init__(self, host_list):
        self.__mc = memcache.Client(host_list)

    def set(self, key, value, time):
        result = self.__mc.set(key, value, time)
        return result

    def get(self, key):
        name = self.__mc.get(key)
        return name

    def delete(self, key):
        result = self.__mc.delete(key)
        return result


def gettoken(corpid, corpsecret, logger):
    gettoken_url = 'https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=' + corpid + '&corpsecret=' + corpsecret
    try:
        token_file = urllib2.urlopen(gettoken_url)
    except urllib2.HTTPError as e:
        print e.code
        print e.read().decode("utf8")

        sys.exit()

    token_data = token_file.read().decode('utf-8')
    token_json = json.loads(token_data)
    token_json.keys()
    token = token_json['access_token']
    logger.info("get access token success, value is:{}".format(token))
    return token


def send_weixin(corp_id, corp_secret, user, content, logger):
    mc = MC(['127.0.0.1:11211'])
    if not mc.get('wx-access-token'):
        access_token = gettoken(corp_id, corp_secret, logger)
        logger.info("key 'wx-access-token' is not in memcached, generate it.")
        mc.set('wx-access-token', access_token, 3600)
    else:
        access_token = mc.get('wx-access-token')
        logger.info("get key 'wx-access-token' from memcached. value is: {}".format(access_token))

    send_url = 'https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=' + access_token
    send_values = {
        "touser": user,  # 企业号中的用户帐号，在zabbix用户Media中配置，如果配置不正常，将按部门发送。
        "toparty": "1",  # 企业号中的部门id
        "msgtype": "text",  # 企业号中的应用id，消息类型。
        "agentid": "1000004",
        "text": {
            "content": content
        },
        "safe": "0"
    }
    send_data = json.dumps(send_values, ensure_ascii=False)
    send_request = urllib2.Request(send_url, send_data)
    response = json.loads(urllib2.urlopen(send_request).read())
    logger.info("send weixin result: {}".format(str(response)))


def send_message(content):
    logger = log("send-weixin.log")
    mc = MC(['127.0.0.1:11211'])
    content_md5 = hashlib.md5(content).hexdigest()
    if not mc.get(content_md5):
		# 账户 token
        send_weixin('', '', '@all', content, logger)
        mc.set(content_md5, content, 1800)
    else:
        logger.info("{} alarm message cached half an hour, don't send it.".format(content))

