package main

import "fmt"

func main() {
	map1 := make(map[string]int)
	map1["one"] = 1
	map1["two"] = 2

	map2 := map[string]int{"one": 1, "two": 2}

	for key, value := range map1 {
		fmt.Println(key, value)
	}

	//检验key 是否存在
	if value, isPresent := map2["one"]; isPresent {
		fmt.Println(value, isPresent)
	}

	//删除key
	delete(map2, "two")
	fmt.Println(map2)

	//map的排序是无序的

}
