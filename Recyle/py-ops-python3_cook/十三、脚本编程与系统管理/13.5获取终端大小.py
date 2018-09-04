# 获取当前终端的大小，以便格式化输出
import os

sz = os.get_terminal_size()
print(sz.columns)
print(sz.lines)
