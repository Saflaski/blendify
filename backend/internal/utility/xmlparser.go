package utility

import "encoding/xml"

type LfmResponse struct {
	XMLName xml.Name `xml:"lfm"`
	Status  string   `xml:"status,attr"`
	Session Session  `xml:"session"`
}

type Session struct {
	Name       string `xml:"name"`
	Key        string `xml:"key"`
	Subscriber int    `xml:"subscriber"`
}

func ParseXMLSessionKey(xmlBody []byte) *LfmResponse {
	var result LfmResponse
	if err := xml.Unmarshal(xmlBody, &result); err != nil {
		panic(err)
	}
	return &result
}
