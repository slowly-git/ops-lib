import random, base64
from hashlib import sha1

data = 's30@@192.168.1.1*4dsff@@192.168.0.205*pycharm@@resin@@@@'
key = 'ctmj#&amp;8hrgow_^sj$ejt@9fzsmh_o)-=(byt5jmg=e3#foya6u'


def crypt(data, key):
    """RC4 algorithm"""
    x = 0

    box = [x for x in range(256)]
    for i in range(256):
        x = (x + box[i] + key[i % len(key)]) % 256
        box[i], box[x] = box[x], box[i]
    x = y = 0
    out = []
    for char in data:
        x = (x + 1) % 256
        y = (y + box[x]) % 256
        box[x], box[y] = box[y], box[x]
        out.append(chr(ord(char) ^ box[(box[x] + box[y]) % 256]))

    return ''.join(out)


def tencode(data, key, encode=base64.b64encode, salt_length=16):
    """RC4 encryption with random salt and final encoding"""
    salt = ''
    for n in range(salt_length):
        salt += chr(random.randrange(256))
    data = salt + crypt(data, sha1((key + salt).encode()).digest())
    if encode:
        data = encode(data.encode())
    return data


def tdecode(data, key, decode=base64.b64decode, salt_length=16):
    """RC4 decryption of encoded data"""
    data = data.decode()
    if decode:
        data = decode(data).decode()
    salt = data[:salt_length]
    return crypt(data[salt_length:], sha1((key + salt).encode()).digest())


aaa = tencode(data, key)
print(aaa)
bbb = tdecode(aaa, key)
print(bbb)

# print(bbb)
