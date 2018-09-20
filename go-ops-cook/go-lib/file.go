package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	//创建文件
	f, err := os.OpenFile("/tmp/test.log", os.O_CREATE|os.O_WRONLY, 0755)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}

	s := "test"

	fmt.Fprintf(f, "测试写文件 %s\n", s)

	f2, err := os.Open("/tmp/test.log")
	defer f2.Close()
	if err != nil {
		fmt.Println(err)
	}
	//用bufio读文件
	s2 := bufio.NewReader(f2)
	//读到\n结束
	str2,err := s2.ReadString('\n')
	if err ==io.EOF{
		fmt.Println("已经读到文件结尾")
	}
	if err != nil{
		fmt.Println(err)
	}

	fmt.Println(str2)
}
