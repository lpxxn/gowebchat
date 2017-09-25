package wxchat

import (
	"crypto/tls"
	"fmt"
	"github.com/lpxxn/gowebchat/utils"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strconv"
	"time"
	"log"
)

type WeChat struct {
	Uuid      string
	TimeStamp string
	DeviceId  string
	Client    *http.Client
	Log *log.Logger
}

func NewWeChat(logger *log.Logger) (*WeChat, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.New("create cookie error")
	}

	transport := *(http.DefaultTransport.(*http.Transport))
	transport.ResponseHeaderTimeout = 2 * time.Minute
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	return &WeChat{
		TimeStamp: strconv.FormatInt(utils.CurrentTimeStep, 10),
		DeviceId:  utils.RandomString("e", 15),
		Client: &http.Client{
			Transport: &transport,
			Jar:       jar,
			Timeout:   transport.ResponseHeaderTimeout,
		},
		Log: logger,
	}, nil
}

/*
	得到UUID
*/
func (weChat *WeChat) GetUuid() error {
	JsLoginUrl := utils.JsLogin + "?appid=" + utils.AppID + "&redirect_uri=" +
		utils.RedirectUri + "&fun=new&lang=zh_CN&_=" + strconv.FormatInt(utils.CurrentTimeStep, 10)

	client := &http.Client{}
	fmt.Println(JsLoginUrl)
	resp, err := client.Get(JsLoginUrl)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error :", err)
		return err
	}
	re := regexp.MustCompile(`window.QRLogin.code = (\d+); window.QRLogin.uuid = "([\s\S]*)";`)
	pm := re.FindStringSubmatch(string(data))
	fmt.Printf("%v \n", pm)
	if len(pm) > 0 {
		code := pm[1]
		if code != "200" {
			readErr := errors.New("the status error")
			fmt.Println(err)
			return readErr
		} else {
			uuid := pm[2]
			fmt.Println("uuid", uuid)
			weChat.Uuid = uuid
		}
	} else {
		err = errors.New("uuid error")
		fmt.Println(err)
		return err
	}
	return nil
}
