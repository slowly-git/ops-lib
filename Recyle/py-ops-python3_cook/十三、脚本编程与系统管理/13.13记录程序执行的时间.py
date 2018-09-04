import time


class Timer:
    """
    定义了一个timer类，可以人为控制动作(star,stop ...)
    """

    def __init__(self, func=time.perf_counter):
        self.elapsed = 0.0
        self._func = func
        self._start = None

    def start(self):
        if self._start is not None:
            raise RuntimeError('Already started')
        self._start = self._func()

    def stop(self):
        if self._start is None:
            raise RuntimeError('Not started')

        end = self._func()
        self.elapsed += end - self._start
        self._start = None

    def reset(self):
        self.elapsed = 0.0

    @property
    def running(self):
        return self._start is not None

    def __enter__(self):
        self.start()
        return self

    def __exit__(self, *args):
        self.stop()


# Timer类用法举例
def countdown(n):
    time.sleep(n)


# 使用一：显式调用 start/stop
t = Timer()
t.start()
countdown(2)
t.stop()
print(t.elapsed)

# 使用二：作为t的上下文管理器
countdown(4)
print(t.elapsed)

with Timer() as t2:
    countdown(1)
print(t2.elapsed)
