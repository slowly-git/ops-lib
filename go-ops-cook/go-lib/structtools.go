package main

import (
	"fmt"
	"math/rand"
	"reflect"
)

type SpotPrice struct {
	Region                     string
	InstenceType               string
	NowPrice                   float64
	LastThreeMonthAveragePrice float64
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
}
