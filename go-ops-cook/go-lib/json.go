package main

import (
	"encoding/json"
	"fmt"
)

type Student struct {
	Name  string `json:"student_name"`
	Age   int `json:"age"`
	Score int `json:"score"`
}

func main() {
	stu := &Student{
		Name:  "stu",
		Age:   12,
		Score: 100,
	}

	//json打包的时候在另外一个包里面，如果结构体成员小写，则不能访问，可以通过tag修改打包后的字段名
	data, err := json.Marshal(stu)
	if err != nil {
		fmt.Println("json faild erra:", err)
	}

	fmt.Println(string(data))
}