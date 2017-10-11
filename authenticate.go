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

// EndUserInfo is information about the enduser that will be sent to bankid-api
type EndUserInfo struct {
	XMLName      xml.Name `xml:"endUserInfo"`
	UserInfoType string   `xml:"type"`
	Value        string   `xml:"value"`
}

// AuthResponse is the response from a bankid auth-request
type AuthResponse struct {
	XMLName        xml.Name `xml:"AuthResponse"`
	OrderRef       string   `xml:"orderRef"`
	AutoStartToken string   `xml:"autoStartToken"`
}

// Authenticate is a method to call the Authenticate resource on the BankID API.
func (c *Client) Authenticate(ssn string, u *EndUserInfo) (*AuthResponse, error) {
	respEnvelope := new(soapEnvelope)
	respEnvelope.Body = soapBody{Content: &AuthResponse{}}

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

	resp, ok := respEnvelope.Body.Content.(*AuthResponse)
	if !ok {
		return nil, errors.New("authResp not ok")
	}
	log.Printf("AutostartToken: %v\n OrderRef: %v\n", resp.AutoStartToken, resp.OrderRef)
	return resp, nil
}
