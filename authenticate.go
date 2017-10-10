package bankid

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
)

type authenticateRequest struct {
	XMLName xml.Name `xml:"typ:AuthenticateRequest"`
	//<!--Optional:-->
	Ssn string `xml:"personalNumber"`
	//<!--0 to 20 repetitions:-->
	UserInfo *EndUserInfo `xml:",omitempty"`
	//<!--Optional:-->
	Alternatives []interface{} `xml:",omitempty"`
}

type EndUserInfo struct {
	XMLName      xml.Name `xml:"endUserInfo"`
	UserInfoType string   `xml:"type"`
	Value        string   `xml:"value"`
}

type authResponse struct {
	XMLName        xml.Name `xml:"AuthResponse"`
	OrderRef       string   `xml:"orderRef"`
	AutoStartToken string   `xml:"autoStartToken"`
}

func (c *Client) Authenticate(ssn string, u *EndUserInfo) (*authResponse, error) {
	respEnvelope := new(SOAPEnvelope)
	respEnvelope.Body = SOAPBody{Content: &authResponse{}}

	a := &authenticateRequest{
		Ssn:      ssn,
		UserInfo: u,
	}

	err := c.RoundTripSoap12("Authenticate", a, respEnvelope)
	if err != nil {
		return nil, err
	}

	if respEnvelope.Body.Fault != nil {
		err := fmt.Errorf("errorcode: %v \n string: %v \n faultStatus: %v \n detailed: %v", respEnvelope.Body.Fault.String, respEnvelope.Body.Fault.Code, respEnvelope.Body.Fault.Detail.FaultStatus, respEnvelope.Body.Fault.Detail.Description)
		return nil, err
	}

	resp, ok := respEnvelope.Body.Content.(*authResponse)
	if !ok {
		return nil, errors.New("authResp not ok")
	}
	log.Printf("AutostartToken: %v\n OrderRef: %v\n", resp.AutoStartToken, resp.OrderRef)
	return resp, nil
}
