package main

import (
	"bytes"
	"fmt"
)

func main() {
	//使用bytes.Buffer拼接字符串
	var b bytes.Buffer
	str1 := "abc"
	str2 := "dev"

	b.WriteString(str1)
	b.WriteString(str2)

	fmt.Println(b.String())
}
