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
	//普通写
	//fmt.Fprint(f, "line 1\nline 2\nline3\n")
	//带缓冲区的写
	fWriter := bufio.NewWriter(f)
	fString := "line 1\nline 2\nline3\n"
	for i := 0; i < len(fString); i++ {
		fWriter.WriteString(fString)
	}
	fWriter.Flush()

	//准备读文件
	f2, err := os.Open("/tmp/test.log")
	defer f2.Close()
	if err != nil {
		os.Exit(1)
	}
	//用bufio读文件
	s2 := bufio.NewReader(f2)
	//读到\n结束
	for {
		str2, err := s2.ReadString('\n')
		if err == io.EOF {
			fmt.Println("已经读到文件结尾")
			break
		}
		if err != nil {
			break
		}
		fmt.Println(str2)
	}
}
