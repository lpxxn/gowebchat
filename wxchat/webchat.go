package wxchat

import (
	"crypto/tls"
	"fmt"
	"github.com/lpxxn/gowebchat/utils"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
	"encoding/xml"
	"encoding/json"
	"io"
	"bytes"
)

type WeChat struct {
	Uuid      string
	TimeStamp string
	DeviceId  string
	Client    *http.Client
	Log       *log.Logger
	// 扫码登录返回数据
	Code        string
	RedirectUri string
	LoginPageModel *WebWxNewLoginPageResponse
	WebWxInitModel *WebWxInit
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
		LoginPageModel: &WebWxNewLoginPageResponse{},
		WebWxInitModel: & WebWxInit{},
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

func (w *WeChat) ScanQrAndLogin() (code string, err error) {
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
	} else {
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

func (w *WeChat) NewLoginPage() (err error) {
	if w.RedirectUri == "" {
		err = errors.New("please get RedirectUri")
		return
	}
	resp, err := w.Client.Get(w.RedirectUri)
	if err != nil {
		return
	}
	defer  resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
	wxnewReq := new(WebWxNewLoginPageResponse)
	err = xml.Unmarshal(data, wxnewReq)
	if err != nil {
		fmt.Println("xml Unmarshal error")
		return
	}
	w.LoginPageModel = wxnewReq

	w.WebWxInitModel.BaseRequest.Sid = wxnewReq.Wxsid
	w.WebWxInitModel.BaseRequest.Skey = wxnewReq.Skey
	w.WebWxInitModel.BaseRequest.Uin = wxnewReq.Wxuin
	w.WebWxInitModel.BaseRequest.DeviceID = w.DeviceId

	baseRequestJson, err := json.Marshal(w.WebWxInitModel)
	fmt.Println("json", string(baseRequestJson), " err :", err)

	initUrl :=  fmt.Sprintf("%s?r=%d&lang=zh_CN&pass_ticket=%s", utils.WebWxInitUrl, utils.MakeTimeStame(), wxnewReq.PassTicket)
	fmt.Println("initUrl :", initUrl)

	initData := w.PostUrl(initUrl, bytes.NewReader(baseRequestJson))
	initModel :=  &WxInitModel{}
	json.Unmarshal(initData, initModel)

	///
	url := "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxgetcontact?r=" + strconv.FormatInt(utils.MakeTimeStame(), 10)
	userlistbytes := w.PostUrl(url, nil)
	userList := &UserList{}
	json.Unmarshal(userlistbytes, userList)

	var findId string
	for _, item := range userList.MemberList {
		if item.RemarkName == "郭雪" {
		//if item.RemarkName == "小胖123" {
			findId = item.UserName
			break;
		}
	}

	sendUrl := fmt.Sprintf("%s?lang=zh_CN&pass_ticket=%s", utils.SendMsgUrl, wxnewReq.PassTicket)
	fmt.Println("sendMsgUrl :", sendUrl)
 	id := strconv.FormatInt(time.Now().UnixNano(), 10)
	fmt.Println(id)
	sendMsg := SendMsgData{
		BaseRequest:w.WebWxInitModel.BaseRequest,
		Msg: MsgModel{
			ClientMsgID: id,
			Content: "test 你好，去屎吧~~~~",
			FromUserName: initModel.User.UserName,
			LocalID: id,
			ToUserName: findId,
			Type: 1,
		},
		Scene: 0,
	}
	fmt.Println("sendMsg:", sendMsg)
	sendJson, err := json.Marshal(sendMsg)
	fmt.Println(sendJson)
	if err != nil {
		fmt.Println("Marshal Send Json error", err.Error())
		return
	}

	if err != nil {
		w.Log.Printf("json.Marshal(%v):%v\n", sendJson, err)
	}

	w.PostUrl(sendUrl, bytes.NewReader(sendJson))
	return
}

func (w * WeChat) PostUrl(url string, body io.Reader) []byte {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := w.Client.Do(req)
	if err != nil {
		fmt.Printf("post url error, Url = %s, error= %v", url, err)
		return nil
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("response data : %s", string(data))
	return data

}

func (w *WeChat) SendMsg()  {

}
