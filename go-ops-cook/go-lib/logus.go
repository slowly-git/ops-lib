package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

var log = logrus.New()

func main() {
	//用法一：直接打印日志
	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	logrus.WithFields(logrus.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	logrus.WithFields(logrus.Fields{
		"omg":    true,
		"number": 100,
	}).Fatal("The ice breaks!")

	//用法二：日记打印到文件
	var log = logrus.New()

	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Fatal("A group of walrus emerges from the ocean")
	file.Close()
}

func init() {
	// 以JSON格式为输出，代替默认的ASCII格式
	//Logging.Formatter = new(logrus.JSONFormatter)
	Logging.Formatter = new(logrus.TextFormatter)
	// 以Stdout为输出，代替默认的stderr
	Logging.Out = os.Stdout
	// 设置日志等级
	Logging.Level = logrus.InfoLevel
	// 删除时间戳
	//Logging.Formatter.(*logrus.TextFormatter).DisableTimestamp = true
	//Logging.Formatter.(*logrus.JSONFormatter).DisableTimestamp = true

	//log.Debug("Useful debugging information.")
	//log.Info("Something noteworthy happened!")
	//log.Warn("You should probably take a look at this.")
	//log.Error("Something failed but I'm not quitting.")
	//// 随后会触发os.Exit(1)
	//log.Fatal("Bye.")
	//// 随后会触发panic()
	//log.Panic("I'm bailing.")
}
