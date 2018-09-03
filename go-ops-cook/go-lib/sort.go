package main

import (
	"fmt"
	"math/rand"
)

/*冒泡排序:跟相邻的元素比较a[0]和a[1]，比较，如果比a[1]大，则交换位置,然后a[1]和a[2]比较，依次冒泡到高位
53142
35142
31542
31452
31425(5冒泡到了最高位)
*/
func bSort(a []int) {
	//每次冒泡要遍历一次a
	for j := 0; j < len(a); j++ {
		//对元素a[j]执行冒泡操作,冒泡到最高位
		for i := 0; i < len(a)-1; i++ {
			if a[i] > a[i+1] {
				a[i], a[i+1] = a[i+1], a[i]
			}
		}
	}
	fmt.Println("冒泡结果:", a)
}

/*选择排序:每次遍历选出一个最大值或者最小值,a[0]和a[1]~a[n]比较，选出最大到值放到a[0]，依次循环...
53142
35142
15342（选出了最小值1）
*/
func sSort(a []int) {
	//遍历原始数列
	for k := range a {
		for i := k + 1; i < len(a); i++ {
			if a[k] > a[i] {
				a[k], a[i] = a[i], a[k]
			}
		}
	}
	fmt.Println("选择排序结果:", a)
}

/*插入排序:假设a[0]是一个有序数列，a[1]~a[n]都是无序数列，从a[1]开始遍历无序数列,
依次从无序数列中选出元素插入到有序数列(类似于插入的元素在有序数列中冒泡操作)
5|3142
35|142
315|42	135|42
...
*/
func iSort(a []int) {
	//遍历无序数列
	for i := 1; i < len(a); i++ {
		//遍历有序数列
		for j := i; j > 0; j-- {
			if a[j] > a[j-1] {
				break
			}
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
	fmt.Println("插入排序结果:", a)
}

/*快速排序:对于一个数组，先定位其中一个元素的位置（左边所有元素比它小，右边所有元素比它大）,
这样就把一个数组拆分成左右两个数组，再对左右对数组进行相同的定位操作,当剩下数组的元素只有1个后，
代表位置确定。
*/
func qSort(a []int, left, right int) {
	//左边数组长度超过总长度后停止排序
	if left > right {
		//fmt.Println("快速排序结果:", a)
		return
	}

	//定位val的位置
	val := a[left]
	k := left
	/*保证a[left]左边的数组比右边的数组小就行,例如57892:
	假设k=0为定位的点，a[k]=5为定位的值
	定位a[k]=5,遍历7892发现a[4]=2满足条件，
	a[k]和a[4]元素互换:27895
	a[4]置换为a[k+1]:27897
	a[k]置换为之前定位的值5:25789
	*/
	for i := left + 1; i <= right; i++ {
		if a[i] < val {
			a[k] = a[i]
			a[i] = a[k+1]
			k++
		}
		fmt.Printf("本次定位点a[%d],定位的值%d \n", left, a[left])
	}

	a[k] = val

	//递归:
	//排序左边的数组
	qSort(a, left, k-1)
	//排序右边的数组
	qSort(a, k+1, right)
}

func main() {
	//冒泡
	a := make([]int, 10)
	for k := range a {
		a[k] = rand.Intn(1000)
	}
	fmt.Println("冒泡排序前:", a)
	bSort(a)

	//选择
	b := make([]int, 10)
	for k := range b {
		b[k] = rand.Intn(1000)
	}
	fmt.Println("选择排序前:", b)
	sSort(b)

	//插入
	c := make([]int, 10)
	for k := range c {
		c[k] = rand.Intn(1000)
	}
	fmt.Println("插入排序前:", c)
	iSort(c)

	//快速排序
	d := make([]int, 10)
	for k := range d {
		d[k] = rand.Intn(1000)
	}
	fmt.Println("快速排序前:", d)
	qSort(d, 0, len(d)-1)
}
