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
	Uin string `json:"Sid"`
	Sid string `json:"Sid"`
	Skey string `json:"Skey"`
	DeviceID string `json:"DeviceID"`
}

type WebWxInit struct {
	BaseRequest BaseRequestModel `json:"BaseRequest"`
}


type SendMsg struct {
	BaseRequest BaseRequestModel `json:"BaseRequest"`
	Msg struct{
		ClientMsgId string `json:"client_msg_id"`

	}
}