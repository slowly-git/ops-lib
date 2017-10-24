# coding:utf-8

import ssl

ssl_cet = ssl.SSLContext.load_verify_locations(cafile='localhost.crt', capath='/etc/pki/tls/certs/')
