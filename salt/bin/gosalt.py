import os, sys
import threading

BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.append(BASE_DIR)

from core import main

if __name__ == '__main__':
    run = threading.Thread(target=main.run)
    run1 = threading.Thread(target=main.run)
    run.start()
    run1.start()
    print('请等待结果',threading.active_count())