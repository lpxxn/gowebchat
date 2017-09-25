package utils

import (
	"time"
	"math/rand"
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
func RandomString(prefix string, n  int) string {
	letterRunes := []rune("123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return prefix + string(b)
}

