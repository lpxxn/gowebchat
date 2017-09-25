package utils

import (
	"math/rand"
	"time"
	"os"
)

/*
	Nano to Milliscond

*/
func MakeTimeStame() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

/*
	随机字符串
	eg. RandomString("e", 15)
*/
func RandomString(prefix string, n int) string {
	letterRunes := []rune("123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return prefix + string(b)
}


func CreateFile(name string, data []byte, isAppend bool) error {
	flag := os.O_CREATE | os.O_WRONLY
	if isAppend {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}

	file, err := os.OpenFile(name, flag, 0666)
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = file.Write(data)
	return err
}