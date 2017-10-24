import shutil
import os

src = "/tmp/pear"
dst = "/home/worker/tmp/pear"

# Copy src to dst. (cp src dst)
shutil.copy(src, dst)
# Copy files, but preserve metadata (cp -p src dst)
shutil.copy2(src, dst)
# Copy directory tree (cp -R src dst)
shutil.copytree(src, dst)
# Move src to dst (mv src dst)
shutil.move(src, dst)

# 如果被复制的是一个软链接，以下方式会只复制软链接而不是其指向的内容
shutil.copy2(src, dst, follow_symlinks=True)

# copytree的时候，可以保留目录中的软链接
shutil.copytree(src, dst, symlinks=True)

# 复制过程中忽略一些文件
shutil.copytree(src, dst, ignore=shutil.ignore_patterns('~', '.pyc'))  # 如果要忽略所有软链接 ignore_dangling_symlinks=True

# 要保证脚本的通用性，建议用os.path 切割文件
import os.path

filename = '/Users/guido/programs/spam.py'
os.path.basename(filename)
>> 'spam.py'
os.path.dirname(filename)
>> '/Users/guido/programs'
os.path.split(filename)
>> ('/Users/guido/programs', 'spam.py')
os.path.join('/new/dir', os.path.basename(filename))
>> '/new/dir/spam.py'
# expanduser 自动将路径中包含的~转换成用户目录
os.path.expanduser('~/guido/programs/spam.py')
>> '/Users/guido/programs/spam.py'

# 错误处理
try:
    shutil.copytree(src, dst)
except shutil.Error as e:
    for src, dst, msg in e.args[0]:
    # src is source name # dst is destination name # msg is error message from exception
    print(dst, src, msg)
