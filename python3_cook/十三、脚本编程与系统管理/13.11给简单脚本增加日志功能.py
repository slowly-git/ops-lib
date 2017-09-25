import logging

hostname = 'pycharm'
item = 'test'
filename = 'test log'
mode = 'test'


def main():
    # 配置日志,level=logging.WARNING表示日志比warning级别更高才会被记录
    logging.basicConfig(
        filename='app.log',
        format='%(levelname)s:%(asctime)s:%(message)s',
        # datefmt='%m/%d/%Y %I:%M:%S %p',
        level=logging.DEBUG
    )
    logging.critical('Host %s unknown', hostname)
    logging.error('Couldn’t find %r', item)
    logging.warning('Feature is deprecated')
    logging.info('Opening file %r, mode=%r', filename, mode)
    logging.debug('Got here')


if __name__ == '__main__':
    main()
