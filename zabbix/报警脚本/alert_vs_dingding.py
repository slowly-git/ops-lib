#!/home/worker/python27/bin/python
# -*-coding: utf-8 -*-

import datetime
import sys
import httplib2
import getopt
import re
import time
import urllib
import smtplib
from email.mime.text import MIMEText
from email.Header import Header
from weixin import send_message
from dingding import send_dingding_msg
from yunpian_python_sdk.model import constant as YC
from yunpian_python_sdk.ypclient import YunpianClient

#smtp 信息, ME USER_NAME都是账户
SMTP_SERVER = 'smtp.mxhichina.com'
ME = 'monitor@tiejin.me'
USER_NAME = 'monitor@tiejin.me'
USER_PASS = 'Closer2018'


nowtime = time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(time.time()))


class SendMessage(object):
    def __init__(self, to_mobiles, message):
        self._to_mobiles = to_mobiles
        self._message = message
        self._apikey = 'f70d99016d7b11573f16f87217b913e2'
        #'【Closer监控】告警:#msg#'
        self._tpl_id = 2280910

    def log(self, message):
        '''日志记录'''
        log = open("/tmp/zabbix_sendsms.log", "a")

        now = datetime.datetime.now()
        current = now.strftime("%Y-%m-%d_%H:%M:%S")

        log.write(current + '\t' + message + '\r\n')
        log.close()

    def sendsms(self):
        """发送短信"""
        now = time.strftime("%Y-%m-%d,%H:%M")
        message = now + ',' + self._message
        tpl_value = urllib.urlencode({"#msg#":message})
        clnt = YunpianClient(self._apikey)
        param = {YC.MOBILE:self._to_mobiles,YC.TPL_ID:self._tpl_id,YC.TPL_VALUE:tpl_value}
        try:
             r = clnt.sms().tpl_single_send(param)
             if r.code() == 0 and r.msg() == "发送成功":
                self.log("message send sucess".decode('gbk').encode('utf8')) 
             else:
                self.log("message send fail:".decode('gbk').encode('utf8')+r.msg())
        except: 
             self.log("message send fail".decode('gbk').encode('utf8'))
        
def msg_args():
    args = sys.argv[1:]
    send_opts = {}

    mobiles = []
    if ',' in args[0]:
        for item in args[0].split(','):
            if not re.search(r'@', item):
                mobiles.append(item)
        send_opts['to_mobiles'] = ','.join(map(str, mobiles))
    else:
        send_opts['to_mobiles'] = args[0]
    send_opts['msg'] = args[1:]
    sendlog = SendMessage(send_opts['to_mobiles'], send_opts['msg'])
    sendlog.log(args[0])

    return send_opts


def send_sms():
    sendargs = msg_args()
    if sendargs['to_mobiles']:
        sendsms = SendMessage(sendargs['to_mobiles'], sendargs['msg'][0].decode('utf-8'))
        sendsms.log(sendargs['to_mobiles'])
        sendsms.log(str(sendargs['msg'][0]).strip('\n'))
        sendsms.sendsms()


def mail_args():
    args = sys.argv[1]
    mails = []
    if ',' in args:
        for item in args.split(','):
            if re.search(r'@', item):
                mails.append(item)

        return ','.join(map(str, mails))
    else:
        return args


def send_mail(to_list, subject, content):
    content = content + '\r\n' + nowtime
    msg = MIMEText(content, 'plain', 'utf-8')
    msg['Subject'] = Header(subject, 'utf-8')
    msg['From'] = ME
    msg['to'] = to_list

    try:
        s = smtplib.SMTP()
        s.connect(SMTP_SERVER)
        s.login(USER_NAME, USER_PASS)
        s.sendmail(ME, to_list, msg.as_string())
        s.close()
        return True
    except Exception, e:
        print str(e)
        return False


def send():
    to_mails = mail_args()
    with open("/tmp/closer_send_msg.log", 'a') as fs:
            fs.write(sys.argv[3])
    if re.search(r'High|Disaster|Warning', sys.argv[3], re.M):
        send_sms()
        try:
            send_message(sys.argv[2])
        except Exception as e:
            with open("/tmp/send.log", 'a') as fd:
                fd.write(str(e))
            pass
        try:
            send_dingding_msg(sys.argv[2])
        except Exception as e:
            with open("/tmp/send_dingding.log", 'a') as fd:
                fd.write(str(e))
            pass
        if to_mails:
            send_mail(to_mails, sys.argv[2], sys.argv[3])
    else:
        if to_mails:
            send_mail(to_mails, sys.argv[2], sys.argv[3])


if __name__ == '__main__':
    send()
