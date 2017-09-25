package wxchat

import (
	"fmt"
	"testing"
)

func TestUuid(t *testing.T) {
	chat, _ := NewWeChat(nil)

	err := chat.GetUuid()
	chat.QrCode()
	fmt.Println("err :", err, "  uuid :", chat.Uuid)

}
