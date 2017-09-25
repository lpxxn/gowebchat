package main

import (
	"github.com/lpxxn/gowebchat/wxchat"
	"os"
	"log"
	"fmt"
	"github.com/lpxxn/gowebchat/utils"
	"path/filepath"
)

func Init() {
	utils.RootPath, _ = os.Getwd()
	utils.ImgPath = filepath.Join(utils.RootPath, "img")
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
