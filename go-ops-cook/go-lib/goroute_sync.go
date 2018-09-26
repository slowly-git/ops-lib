package main

import "fmt"

/*
用chan同步的方式，求10W以内的素数;
管道的关闭要在管道停止使用后;
for range chan会便利chan,且检测管道是否关闭;
*/

func main() {
	intChan := make(chan int, 1000)
	resChan := make(chan int, 1000)
	exitChan := make(chan bool, 8)

	//生产10W以内的数
	go func() {
		for i := 0; i < 100000; i++ {
			intChan <- i
		}

		//在管道中存入10W个数后关闭管道
		close(intChan)
	}()

	//启动8个goroute去计算
	for i := 0; i < 8; i++ {
		go calc(intChan, resChan, exitChan)
	}

	//等待4个goroute全部退出
	go func() {
		for i := 0; i < 8; i++ {
			<-exitChan //取出来的值直接扔掉
			fmt.Println("wait goroute ", i, "exited!")
		}

		//启用的8个goroute都退出后，关闭管道resChan
		close(resChan)
		close(exitChan)
	}()

	//取素数结果
	for v := range resChan {
		fmt.Println(v)
		//_ = v
	}
}

func calc(task, result chan int, exitChan chan bool) {
	for v := range task {

		flag := true
		//素数判断
		for i := 2; i < v; i++ {
			if (v % i) == 0 {
				flag = false
				break
			}
		}

		if flag {
			result <- v
		}
	}

	//传递完成信号给关闭Chan
	exitChan <- true
}
