package bankid

import (
	"encoding/xml"
	"errors"
)

type fileSignRequest struct {
	XMLName xml.Name `xml:"typ:FileSignRequest"`
	//<!--Optional:-->
	Ssn         string `xml:"personalNumber"`
	VisibleData string `xml:"userVisibleData"`
	//<!--Optional:-->
	NonVisibleData string `xml:"userNonVisibleData"`
	FileName       string `xml:"fileName"`
	FileContent    string `xml:"fileContent"`
	//<!--0 to 20 repetitions:-->
	UserInfo *EndUserInfo `xml:",omitempty"`
	//<!--Optional:-->
	Alternatives []interface{} `xml:",omitempty"`
}

func (c *Client) FileSign() (string, error) {
	return "", errors.New("Not implemented")
}
