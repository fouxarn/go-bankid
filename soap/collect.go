package soap

import (
	"encoding/xml"
	"fmt"
)

type orderRef struct {
	XMLName  xml.Name `xml:"typ:orderRef"`
	OrderRef string   `xml:",chardata"`
}

type collectResponse struct {
	XMLName   xml.Name       `xml:"CollectResponse"`
	Status    progressStatus `xml:"progressStatus"`
	Signature string         `xml:"signature"`
	UserInfo  userInfo       `xml:"userInfo"`
}

type progressStatus string

const (
	StatusOutstandingTransaction progressStatus = "OUTSTANDING_TRANSACTION"
	StatusNoClient               progressStatus = "NO_CLIENT"
	StatusStarted                progressStatus = "STARTED"
	StatusUserSign               progressStatus = "USER_SIGN"
	StatusUserReq                progressStatus = "USER_REQ"
	StatusComplete               progressStatus = "COMPLETE"
)

type userInfo struct {
	GivenName      string `xml:"givenName"`
	Surname        string `xml:"surname"`
	Name           string `xml:"name"`
	PersonalNumber string `xml:"personalNumber"`
	NotBefore      string `xml:"notBefore"`
	NotAfter       string `xml:"notAfter"`
	IPAddress      string `xml:"ipAddress"`
}

func (c *Client) Collect(ref string) (*collectResponse, error) {
	respEnvelope := new(SOAPEnvelope)
	respEnvelope.Body = SOAPBody{Content: &authResponse{}}

	orderRef := &orderRef{
		OrderRef: ref,
	}

	respEnvelope.Body = SOAPBody{Content: &collectResponse{}}

	err := c.RoundTripSoap12("Collect", orderRef, respEnvelope)
	if err != nil {
		return nil, err
	}

	collResp, ok := respEnvelope.Body.Content.(*collectResponse)
	if ok {
		return collResp, nil
	}
	return nil, fmt.Errorf("errorcode: %v \n string: %v \n faultStatus: %v \n detailed: %v", respEnvelope.Body.Fault.String, respEnvelope.Body.Fault.Code, respEnvelope.Body.Fault.Detail.FaultStatus, respEnvelope.Body.Fault.Detail.Description)
}
