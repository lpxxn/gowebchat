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
