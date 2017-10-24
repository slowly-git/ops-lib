from core import salt_api
from conf import settings


def run():
    salt = salt_api.SaltApi(**settings.SALT_INFO)
    auth_token = salt.uer_login()
    minions = salt.list_all_key()
    print(auth_token)
    for x in minions:
        print(x)


def gethost():
    pass
