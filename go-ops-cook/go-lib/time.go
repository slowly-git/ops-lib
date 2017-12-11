package main

import (
	"time"
	"fmt"
	"reflect"
)

func main() {
	// get now time UTC
	tNowUtc := time.Now().UTC()
	// get now time
	tNow := time.Now()

	// time to string , can't change layout
	timeNow := tNow.Format("2006-01-02 15:04:05")
	fmt.Println(timeNow)
	fmt.Println(reflect.TypeOf(timeNow))

	// string to time , layout must be "2006-01-02 15:04:05"
	t, _ := time.Parse("2006-01-02 15:04:05", "2017-01-06 08:09:10")
	fmt.Println(t)

	// 计算时间
	var week time.Duration
	week = 60 * 60 * 24 * 7 * 1e9 // must be in nanosec
	week_from_now := tNow.Add(week)
	fmt.Println(week_from_now)

	// other format
	fmt.Println(tNowUtc.Format("02 Jan 2006 15:04"))
	fmt.Println(tNowUtc.Format("20060102"))
	fmt.Println(tNowUtc.Format(time.RFC822))
	fmt.Println(tNowUtc.Format(time.ANSIC))

}
