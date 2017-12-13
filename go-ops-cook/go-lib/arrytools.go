package main

import (
	"log"
	"sort"
)

func main() {
	//先排序再去重，
	array := []int{1, 3, 3, 3, 2, 1, 2, 1}

	sort.Ints(array)

	i, j := 1, 0
	for i < len(array) {
		if array[i] != array[j] {
			j++
			array[i], array[j] = array[j], array[i]
		}
		i++
	}
	array = array[0: j+1]

	log.Print(array)

	//用map去重，数据位置不变
	arraymap := []int{1, 3, 3, 3, 2, 1, 2, 1}

	var result []int
	dic := make(map[int]bool)

	for _, element := range arraymap {
		if !dic[element] {
			dic[element] = true
			result = append(result, element)
		}
	}

	log.Print(result)

}
