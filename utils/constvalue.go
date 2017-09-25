package utils

import (
)

const (
	AppID       = "wx782c26e4c19acffb"
	JsLogin     = "https://login.wx.qq.com/jslogin"
	RedirectUri = "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage"
	UserAgent   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"
	Lang        = "zh_CN"
)

var (
	CurrentTimeStep = MakeTimeStame()
)

var Root string
