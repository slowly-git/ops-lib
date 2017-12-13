package main

import (
	"flag"
	"fmt"
)

func main() {
	name := flag.String("name", "", "what is your name?")
	age := flag.Int("age", 22, "how old aer you?")
	//命令行参数必须使用" -married=true|false"
	married := flag.Bool("married", false, "are you marride? married = false|true")

	//直接使用变量地址的写法
	var address string
	flag.StringVar(&address, "address", "SiChuan", "where is you address?")

	//解析命令行输出的参数
	flag.Parse()

	fmt.Println("参数 name:", *name)
	fmt.Println("参数 age:", *age)
	fmt.Println("参数 married:", *married)
	fmt.Println("参数 address:", address)
}
