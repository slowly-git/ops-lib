package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	//使用bytes.Buffer拼接字符串,效率最高
	var b bytes.Buffer
	str1 := "abc"
	str2 := "dev"

	b.WriteString(str1)
	b.WriteString(str2)

	fmt.Println(b.String())

	//使用Sprintf拼接
	hello := "hello"
	world := "world"

	fmt.Sprintf("%s,%s", hello, world)

	//使用，json拼接,适合拼接字符串数组
	for i := 0; i < 100; i++ {
		strings.Join([]string{hello, world}, ",")
	}

}
