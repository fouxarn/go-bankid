package bankid

import (
	"encoding/xml"
	"errors"
)

type signRequest struct {
	XMLName xml.Name `xml:"typ:SignRequest"`
	//<!--Optional:-->
	Ssn         string `xml:"personalNumber"`
	VisibleData string `xml:"userVisibleData"`
	//<!--Optional:-->
	NonVisibleData string `xml:"userNonVisibleData"`
	//<!--0 to 20 repetitions:-->
	UserInfo *EndUserInfo `xml:",omitempty"`
	//<!--Optional:-->
	Alternatives []interface{} `xml:",omitempty"`
}

// Sign is a method to call the Sign resource on the BankID API.
func (c *Client) Sign() (string, error) {
	return "", errors.New("Not implemented")
}
