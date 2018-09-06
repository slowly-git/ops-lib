package main

import (
	"fmt"
	"reflect"
	"math/rand"
)

type SpotPrice struct {
	Region                     string
	InstenceType               string
	NowPrice                   float64
	LastThreeMonthAveragePrice float64
}

//定义一个链表结构,单链表
type Student struct {
	name  string
	age   int
	score float32
	next  *Student
}

//遍历链表
func trans(p *Student) {
	for p != nil {
		fmt.Println(*p)
		p = p.next
	}
}

/*链表：尾部插入法*/
func insertTrail(p *Student) {
	var tail = p
	for i := 0; i < 10; i++ {
		var stu = &Student{
			name:  fmt.Sprintf("stu%d", i),
			age:   rand.Intn(100),
			score: rand.Float32() * 100,
		}
		tail.next = stu
		tail = stu
	}

	trans(p)
}

/*链表：头部插入法*/
func insertHead(p *Student) {
	for i := 0; i < 10; i++ {
		var stu = &Student{
			name:  fmt.Sprintf("stu%d", i),
			age:   rand.Intn(100),
			score: rand.Float32() * 100,
		}
		stu.next = p
		p = stu
	}

	trans(p)
}

func main() {

	s := SpotPrice{"ap", "c4.large", 0.22, 0.99}
	str := []string{"a", "b", "c"}

	//struct
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	n := t.NumField()

	fmt.Println(t)
	fmt.Println(v, v.Kind())
	fmt.Println(n)

	for i := 0; i < n; i++ {
		fmt.Println(t.Field(i).Name, t.Field(i).Type)
	}

	//[]string
	sa, sp := reflect.ValueOf(str), reflect.ValueOf(&str).Elem()
	//st := reflect.TypeOf(str)

	fmt.Println(sa.CanAddr(), sa.CanSet())
	fmt.Println(sp.CanAddr(), sp.CanSet())

	sp.Index(0).SetString("xxxx")

	fmt.Println(str)

	//创建普通二维数组
	arr2Dim := [2][3]int{{1, 2, 3}, {4, 5, 6}}
	for i := range arr2Dim {
		for j := range arr2Dim[i] {
			fmt.Printf("%v ", arr2Dim[i][j])
		}
		fmt.Println()
	}

	//创建二维Slice
	array := make([][]int, 5)
	for i := range array {
		subArray := make([]int, i+1)
		for j := range subArray {
			subArray[j] = j + 1
		}
		array[i] = subArray
	}
	// 输出
	for i := range array {
		for j := range array[i] {
			fmt.Printf("%v ", array[i][j])
		}
		fmt.Println()

	}

	/*链表*/
	//定义头节点
	var head Student
	head.name = "head"
	head.age = 18
	head.score = 100

	//insertTrail(&head)
	insertHead(&head)
}
