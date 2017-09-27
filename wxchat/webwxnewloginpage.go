package wxchat

import "encoding/xml"

type WebWxNewLoginPageResponse struct {
	XmlName     xml.Name `xml:"xml_name"`
	Ret         int      `xml:"ret"`
	Message     string   `xml:"message"`
	Skey        string   `xml:"skey"`
	Wxsid       string   `xml:"wxsid"`
	Wxuin       string   `xml:"wxuin"`
	PassTicket  string   `xml:"pass_ticket"`
	Isgrayscale int      `xml:"isgrayscale"`
}

type BaseRequestModel struct {
	Uin      string `json:"Sid"`
	Sid      string `json:"Sid"`
	Skey     string `json:"Skey"`
	DeviceID string `json:"DeviceID"`
}

type WebWxInit struct {
	BaseRequest BaseRequestModel `json:"BaseRequest"`
}

type MsgModel struct {
	ClientMsgID  string  `json:"ClientMsgId"`
	Content      string `json:"Content"`
	FromUserName string `json:"FromUserName"`
	LocalID      string  `json:"LocalID"`
	ToUserName   string `json:"ToUserName"`
	Type         int    `json:"Type"`
}

type SendMsgData struct {
	BaseRequest BaseRequestModel `json:"BaseRequest"`
	Msg         MsgModel         `json:"Msg"`
	Scene       int64            `json:"Scene"`
}


type User struct {
	UserName          string `json:"UserName"`
	Uin               int64  `json:"Uin"`
	NickName          string `json:"NickName"`
	HeadImgUrl        string `json:"HeadImgUrl" xml:""`
	RemarkName        string `json:"RemarkName" xml:""`
	PYInitial         string `json:"PYInitial" xml:""`
	PYQuanPin         string `json:"PYQuanPin" xml:""`
	RemarkPYInitial   string `json:"RemarkPYInitial" xml:""`
	RemarkPYQuanPin   string `json:"RemarkPYQuanPin" xml:""`
	HideInputBarFlag  int    `json:"HideInputBarFlag" xml:""`
	StarFriend        int    `json:"StarFriend" xml:""`
	Sex               int    `json:"Sex" xml:""`
	Signature         string `json:"Signature" xml:""`
	AppAccountFlag    int    `json:"AppAccountFlag" xml:""`
	VerifyFlag        int    `json:"VerifyFlag" xml:""`
	ContactFlag       int    `json:"ContactFlag" xml:""`
	WebWxPluginSwitch int    `json:"WebWxPluginSwitch" xml:""`
	HeadImgFlag       int    `json:"HeadImgFlag" xml:""`
	SnsFlag           int    `json:"SnsFlag" xml:""`
}


type WxInitModel struct {
	User                User    `json:"User"`
	ContactList         []User  `json:"ContactList"`
}

type UserList struct {
	MemberList []User `json:"MemberList"`
}