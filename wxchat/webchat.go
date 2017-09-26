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
	"path/filepath"
)

type WeChat struct {
	Uuid      string
	TimeStamp string
	DeviceId  string
	Client    *http.Client
	Log *log.Logger
	// 扫码登录返回数据
	Code string
	RedirectUri string
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

/*
	得到QR图片
 */
func (weChat *WeChat) QrCode() error {
	if weChat.Uuid == "" {
		return errors.New("Uuid is empty")
	}
	utils.RemoveAllInDir(utils.ImgPath)
	var qrUrl = utils.LoginQRImg + weChat.Uuid

	req, err := http.NewRequest("GET", qrUrl, nil)
	if err != nil {
		weChat.Log.Fatalln(err)
		return err
	}

	resp, err := weChat.Client.Do(req)
	if err != nil {
		weChat.Log.Fatalln(err)
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		weChat.Log.Fatalln(err)
		return err
	}
	name := utils.RandomString("qr", 5) + ".jpg"
	pathb := filepath.Join(utils.ImgPath, name)
	fmt.Println(pathb)
	return utils.CreateFile(pathb, data, true)
}

func (w *WeChat) ScanQrAndLogin() (code string, err error){
	timeStep := strconv.FormatInt(utils.MakeTimeStame(), 10)
	loginUrl := fmt.Sprintf("%s?loginicon=false&uuid=%s&tip=0&_=%s", utils.ScanORLogin, w.Uuid, timeStep)
	fmt.Println(loginUrl)
	resp, err := w.Client.Get(loginUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	strData := string(data)
	regCode := regexp.MustCompile(`window.code=(\d+)`)
	fss := regCode.FindStringSubmatch(strData)
	if len(fss) > 0 {
		code = fss[1]
	} else{
		err = errors.New("no code")
		return
	}

	switch code {
	case "201":
		fmt.Println("login use phone")
	case "200":
		regUrl := regexp.MustCompile(`window.redirect_uri="([\s\S]*)"`)
		fssUrl := regUrl.FindStringSubmatch(strData)
		if len(fssUrl) > 0 {
			w.RedirectUri = fmt.Sprintf("%s&fun=new", fssUrl[1])
			fmt.Println(w.RedirectUri)
		} else {
			err = errors.New("can not find Redirect Uil")
			return
		}
	case "408":
		fmt.Println("time out rerequest")
	default:
		err = errors.New("unknown error cede is " + code)

	}
	return
}