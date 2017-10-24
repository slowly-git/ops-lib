# 查找文件，可使用 os.walk() 函数，传一个顶级目录名给它。
# 下面是一个例子，查找特定的文件名并答应所有符合条件的文件全路径：
import os
import sys


def findfile(start, name):
    for relpath, dirs, files in os.walk(start):
        if name in files:
            full_path = os.path.join(start, relpath, name)
            print(os.path.normpath(os.path.abspath(full_path)))


if __name__ == '__main__':
    findfile(sys.argv[1], sys.argv[2])

# 例二:找出近期被修改过的文件
import time
import os, sys


def modified_within(top, seconds):
    now = time.time()
    for path, dirs, files in os.walk(top):
        for name in files:
            fullpath = os.path.join(path, name)
            if os.path.exists(fullpath):
                mtime = os.path.getmtime(fullpath)
                if mtime > (now - seconds):
                    print(fullpath)


if __name__ == '__main__':
    if len(sys.argv) != 3:
        print("Usage: {} dir sechonds".format(sys.argv[0]))
        raise SystemExit(1)
    modified_within(sys.argv[1], float(sys.argv[2]))
