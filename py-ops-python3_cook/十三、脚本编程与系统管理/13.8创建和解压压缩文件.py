import shutil

# 解压文件Python-3.3.0.tgz
shutil.unpack_archive('Python-3.3.0.tgz','/tmp')

# 用zip压缩文件Python-3.3.0.tgz,名字为py33.zip
shutil.make_archive('py33','zip','Python-3.3.0')

# 可以压缩的模式[('bztar', "bzip2'ed tar-file"), ('gztar', "gzip'ed tar-file"),('tar', 'uncompressed tar file'), ('zip', 'ZIP file')]
shutil.get_archive_formats()