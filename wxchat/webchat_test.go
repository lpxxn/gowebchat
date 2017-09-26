package wxchat

import (
	"fmt"
	"testing"
	"os"
	"path/filepath"
	"github.com/lpxxn/gowebchat/utils"
)

func TestUuid(t *testing.T) {
	utils.RootPath, _ = os.Getwd()
	utils.RootPath = filepath.Join(utils.RootPath, "../")
	utils.ImgPath = filepath.Join(utils.RootPath, "img")

	utils.RemoveAllInDir(utils.ImgPath)

	chat, _ := NewWeChat(nil)

	err := chat.GetUuid()
	chat.QrCode()
	fmt.Println("err :", err, "  uuid :", chat.Uuid)

}
