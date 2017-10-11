package bankid

import (
	"encoding/xml"
	"fmt"
)

type orderRef struct {
	XMLName  xml.Name `xml:"typ:orderRef"`
	OrderRef string   `xml:",chardata"`
}

// CollectResponse is the response from a bankid collect-request
type CollectResponse struct {
	XMLName   xml.Name       `xml:"CollectResponse"`
	Status    progressStatus `xml:"progressStatus"`
	Signature string         `xml:"signature"`
	UserInfo  UserInfo       `xml:"userInfo"`
}

type progressStatus string

// Different statuses received from bankid-api
const (
	StatusOutstandingTransaction progressStatus = "OUTSTANDING_TRANSACTION"
	StatusNoClient               progressStatus = "NO_CLIENT"
	StatusStarted                progressStatus = "STARTED"
	StatusUserSign               progressStatus = "USER_SIGN"
	StatusUserReq                progressStatus = "USER_REQ"
	StatusComplete               progressStatus = "COMPLETE"
)

// UserInfo is all information about a user returned from bankid-api
type UserInfo struct {
	GivenName      string `xml:"givenName"`
	Surname        string `xml:"surname"`
	Name           string `xml:"name"`
	PersonalNumber string `xml:"personalNumber"`
	NotBefore      string `xml:"notBefore"`
	NotAfter       string `xml:"notAfter"`
	IPAddress      string `xml:"ipAddress"`
}

// Collect is a method to call the Collect resource on the BankID API.
func (c *Client) Collect(ref string) (*CollectResponse, error) {
	respEnvelope := new(soapEnvelope)
	respEnvelope.Body = soapBody{Content: &AuthResponse{}}

	orderRef := &orderRef{
		OrderRef: ref,
	}

	respEnvelope.Body = soapBody{Content: &CollectResponse{}}

	err := c.RoundTripSoap12("Collect", orderRef, respEnvelope)
	if err != nil {
		return nil, err
	}

	collResp, ok := respEnvelope.Body.Content.(*CollectResponse)
	if !ok {
		err = fmt.Errorf("errorcode: %v \n string: %v \n faultStatus: %v \n detailed: %v", respEnvelope.Body.Fault.String, respEnvelope.Body.Fault.Code, respEnvelope.Body.Fault.Detail.FaultStatus, respEnvelope.Body.Fault.Detail.Description)
		return nil, err
	}
	return collResp, nil
}
