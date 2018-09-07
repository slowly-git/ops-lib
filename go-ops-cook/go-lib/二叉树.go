package main

import "fmt"

type Student struct {
	Name  string
	Age   int
	Score int
	left  *Student
	right *Student
}

func main() {
	left := &Student{
		Name:  "left",
		Age:   11,
		Score: 88,
	}
	right := &Student{
		Name:  "right",
		Age:   33,
		Score: 90,
	}
	root := &Student{
		Name:  "test",
		Age:   22,
		Score: 100,
		left:  left,
		right: right,
	}

	trans(root)
}

func trans(root *Student) {
	if root == nil {
		return
	}
	//前序遍历
	fmt.Println("前序遍历:",root)
	trans(root.left)
	trans(root.right)
	//中序遍历
	//trans(root.left)
	//fmt.Println("中序遍历",root)
	//trans(root.right)
	//后序遍历
	//trans(root.left)
	//trans(root.right)
	//fmt.Println("后序遍历",root)
}
