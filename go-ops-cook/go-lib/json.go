package main

import (
	"encoding/json"
	"fmt"
)

type Student struct {
	Name   string `json:"userName""`
	Age    int
	Family map[int]string
}

func testStruct() string {
	familyMap := make(map[int]string)
	familyMap[0] = "mam"
	familyMap[1] = "bab"

	s := &Student{
		Name:   "stu1",
		Age:    12,
		Family: familyMap,
	}
	sJson, _ := json.Marshal(s)
	fmt.Println("test struct json :", string(sJson))
	return string(sJson)
}

func testMap() {
	m := make(map[string]interface{})
	m["name"] = "stu2"
	m["age"] = 11

	mJson, _ := json.Marshal(m)
	fmt.Println("test map json :", string(mJson))
}

func testSlice() {
	s := make([]map[string]interface{}, 0)

	m1 := make(map[string]interface{})
	m1["name"] = "stu3"
	m1["age"] = 12

	m2 := make(map[string]interface{})
	m2["name"] = "stu4"
	m2["age"] = 18

	s = append(s, m1)
	s = append(s, m2)

	sliceJson, _ := json.Marshal(s)
	fmt.Println("stest Slice json :", string(sliceJson))
}

func main() {
	//structJson := testStruct()
	testMap()
	testSlice()
	testUnMarshal()
}

//反序列化
func testUnMarshal()  {
	date := testStruct()

	var stu Student

	//此处必须传递指针
	err :=json.Unmarshal([]byte(date),&stu)
	if err!=nil{
		fmt.Println(err)
	}

	fmt.Println("反序列化测试:",stu)
}