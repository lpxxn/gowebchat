package main

import (
	"github.com/lpxxn/gowebchat/wxchat"
	"os"
	"log"
	"fmt"
	"github.com/lpxxn/gowebchat/utils"
)

func Init() {
	utils.Root, _ = os.Getwd()
}

func main() {
	fileName := "log.txt"
	logFile, err := os.OpenFile(fileName, os.O_CREATE | os.O_APPEND | os.O_RDWR, 06666)
	defer logFile.Close()
	if err != nil {
		fmt.Println("open file error")
		return
	}

	logger := log.New(logFile, "",log.LstdFlags)
	wechat, _ := wxchat.NewWeChat(logger)

	fmt.Println(wechat.TimeStamp, " device: ", wechat.DeviceId)
	wechat.Log.Fatal("error")
}
