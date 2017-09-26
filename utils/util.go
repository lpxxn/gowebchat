package utils

import (
	"math/rand"
	"time"
	"os"
	"path/filepath"
)

func Init() {

}
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
	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return prefix + string(b)
}

/*
	create a file
 */
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

/*
	remover all contents of a directory
 */
func RemoveAllInDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names , err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}