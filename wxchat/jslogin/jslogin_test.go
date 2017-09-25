package jslogin

import (
	"testing"
	"github.com/lpxxn/gowebchat/wxchat"
	"strconv"
	"net/http"
	"fmt"
	"io/ioutil"
	"regexp"
	"github.com/pkg/errors"
)

var JsLoginUrl = wxchat.JsLogin + "?appid=" + wxchat.AppID + "&redirect_uri=" +
	wxchat.RedirectUri + "&fun=new&lang=zh_CN&_=" + strconv.FormatInt(wxchat.CurrentTimeStep, 10)

func TestUuid(t *testing.T) {
	client := &http.Client{}
	resp, err := client.Get(JsLoginUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error :", err)
	}
	re := regexp.MustCompile(`window.QRLogin.code = (\d+); window.QRLogin.uuid = "([\s\S]*)";`)
	pm := re.FindStringSubmatch(string(data))
	fmt.Printf("%v \n", pm)
	if len(pm) > 0 {
		code := pm[1]
		if code != "200" {
			err := errors.New("the status error")
			fmt.Println(err)
		} else {
			uuid := pm[2]
			fmt.Println("uuid", uuid)
		}
	} else {
		err = errors.New("uuid error")
		fmt.Println(err)
	}

}
