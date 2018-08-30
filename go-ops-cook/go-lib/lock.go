package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

//互斥锁(适用于写多读少,读写都需要加锁)
var lock sync.Mutex

//读写锁(适用于读多写少)
var rwlock sync.RWMutex

//性能对比
var rwCount int32
var muCount int32

func testLock() {
	a := make(map[int]int, 100)
	a[1] = 10
	a[2] = 10
	a[3] = 10
	a[4] = 10

	//模拟互斥锁写
	for i := 0; i < 2; i++ {
		go func(b map[int]int) {
			//加互斥锁
			lock.Lock()
			b[1] = rand.Intn(100)
			//模拟写 10ms
			time.Sleep(time.Millisecond * 10)
			//解锁
			lock.Unlock()
		}(a)
	}

	//模拟读写锁读
	for i := 0; i < 100; i++ {
		go func(b map[int]int) {
			for {
				//加读写锁
				rwlock.RLock()
				//模拟读 1ms
				time.Sleep(time.Millisecond)
				//解锁
				rwlock.RUnlock()

				//利用原子操作计数:原子操作是串行的.
				atomic.AddInt32(&rwCount, 1)
			}
		}(a)
	}
	time.Sleep(time.Second * 3)
	fmt.Printf("读写锁，3秒运行 %d 次 \n",atomic.LoadInt32(&rwCount))

	//模拟互斥锁读
	for i := 0; i < 100; i++ {
		go func(b map[int]int) {
			for {
				//加互斥锁
				lock.Lock()
				//模拟读 1ms
				time.Sleep(time.Millisecond)
				//解锁
				lock.Unlock()

				//利用原子操作计数:原子操作是串行的.
				atomic.AddInt32(&muCount, 1)
			}
		}(a)
	}
	time.Sleep(time.Second * 3)
	fmt.Printf("互斥锁，3秒运行 %d 次",atomic.LoadInt32(&muCount))
}

func main() {
	testLock()
}
