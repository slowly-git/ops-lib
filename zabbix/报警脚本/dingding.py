#!/home/worker/python27/bin/python
# -*- coding: utf-8 -*-
 
import requests
import json
import sys
import os
import hashlib
import memcache
 
headers = {'Content-Type': 'application/json;charset=utf-8'}
api_url = "https://oapi.dingtalk.com/robot/send?access_token=1f0f72b6f6b7e48efd74821e0455606303a9d44c1694ad0f86601262d07a64d3"

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

def send_dingding_msg(text):
    mc = MC(['127.0.0.1:11211'])
    content_md5 = hashlib.md5(text).hexdigest()
    if not mc.get(content_md5):
        json_text= {
         "msgtype": "text",
            "text": {
                "content": text
            },
            "at": {
                "atMobiles": [
                    "xxxxxx"
                ],
                "isAtAll": False
            }
        }
        print requests.post(api_url,json.dumps(json_text),headers=headers).content
        mc.set(content_md5, text, 1800)
     
if __name__ == '__main__':
    text = sys.argv[1]
    send_dingding_msg(text)
