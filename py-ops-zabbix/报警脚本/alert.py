#!/usr/bin/env python
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

#smtp 信息, ME USER_NAME都是账户
SMTP_SERVER = ''
ME = ''
USER_NAME = ''
USER_PASS = ''


nowtime = time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(time.time()))


class SendMessage(object):
    def __init__(self, to_mobiles, message):
        self._to_mobiles = to_mobiles
        self._message = message
        self._send_user = 'pinguofetion'
        self._send_password = 'nPhqjgZbOKi'
        self._url = 'http://service.winic.org:8009/sys_port/gateway/?id=%s&pwd=%s&to=%s&content=%s&time='

        self._send = httplib2.Http()

    def log(self, message):
        '''日志记录'''
        log = open("/dev/shm/zabbix_sendsms.log", "a")

        now = datetime.datetime.now()
        current = now.strftime("%Y-%m-%d_%H:%M:%S")

        log.write(current + '\t' + message + '\r\n')
        log.close()

    def sendsms(self):
        """发送短信"""
        now = time.strftime("%Y-%m-%d,%H:%M")
        message = now + ',' + self._message
        url = self._url % (
        urllib.quote(self._send_user), urllib.quote(self._send_password), urllib.quote(self._to_mobiles),
        urllib.quote(message.encode('gb2312')))
        try:
            response, content = self._send.request(url)
            send_status = content.split("/")
            if response["status"] == "200" and send_status[0] == "000":
                self.log("message send sucess".decode('gbk').encode('utf8'))
            else:
                self.log("message send fail".decode('gbk').encode('utf8'))
        except:
            self.log('request url fail'.decode('gbk').encode('utf8'))


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
    with open("/tmp/xxxxxxxxx.log", 'a') as fs:
	    fs.write(sys.argv[3])
    if re.search(r'High|Disaster|Warning', sys.argv[3], re.M):
        send_sms()
	try:
	    send_message(sys.argv[2])
	except Exception as e:
	    with open("/tmp/send.log", 'a') as fd:
		fd.write(str(e))
	    pass
        if to_mails:
            send_mail(to_mails, sys.argv[2], sys.argv[3])
    else:
        if to_mails:
            send_mail(to_mails, sys.argv[2], sys.argv[3])


if __name__ == '__main__':
    send()
