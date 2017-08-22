package main

import (
	"crypto/tls"
	"encoding/xml"
	"gobankid/soap"
	"log"
	"net/http"
)

type authenticateRequest struct {
	XMLName xml.Name `xml:"typ:AuthenticateRequest"`
	//<!--Optional:-->
	Ssn string `xml:"personalNumber"`
	//<!--0 to 20 repetitions:-->
	UserInfo *endUserInfo `xml:",omitempty"`
	//<!--Optional:-->
	Alternatives []interface{} `xml:",omitempty"`
}

type endUserInfo struct {
	XMLName      xml.Name `xml:"endUserInfo"`
	UserInfoType string   `xml:"type"`
	Value        string   `xml:"value"`
}

type orderRef struct {
	OrderRef string `xml:"typ:orderRef"`
}

type authResponse struct {
	XMLName        xml.Name `xml:"AuthResponse"`
	OrderRef       string   `xml:"orderRef"`
	AutoStartToken string   `xml:"autoStartToken"`
}

func main() {
	cert, err := tls.LoadX509KeyPair("cert.crt", "key.key")
	if err != nil {
		log.Fatal(err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{cert},
		},
	}

	client := &http.Client{Transport: tr}

	//s := soap.NewSOAPClient("https://appapi.test.bankid.com/rp/v4?wsdl", true, nil, cert)
	s := soap.Client{
		URL:       "https://appapi.test.bankid.com/rp/v4?wsdl",
		Namespace: "http://bankid.com/RpService/v4.0.0/types/",
		Config:    client,
	}

	respEnvelope := new(soap.SOAPEnvelope)
	respEnvelope.Body = soap.SOAPBody{Content: &authResponse{}}

	u := &endUserInfo{
		UserInfoType: "IP_ADDR",
		Value:        "192.168.0.1",
	}
	a := &authenticateRequest{
		Ssn:      "190102030400",
		UserInfo: u,
	}
	//err = s.Call("Authenticate", a, t)
	err = s.RoundTripSoap12("Authenticate", a, respEnvelope)
	if err != nil {
		log.Println(err)
	}

	if respEnvelope.Body.Fault != nil {
		log.Printf("errorcode: %v \n string: %v \n faultStatus: %v \n detailed: %v \n", respEnvelope.Body.Fault.String, respEnvelope.Body.Fault.Code, respEnvelope.Body.Fault.Detail.FaultStatus, respEnvelope.Body.Fault.Detail.Description)
		return
	}

	resp, ok := respEnvelope.Body.Content.(*authResponse)
	if ok {
		log.Printf("AutostartToken: %v\n OrderRef: %v\n", resp.AutoStartToken, resp.OrderRef)
	}

	/*orderRef := &orderRef{
		OrderRef: resp.OrderRef,
	}

	respEnvelope.Body = soap.SOAPBody{Content: &collectResponse{}}
	err = s.RoundTripSoap12("Collect", orderRef, respEnvelope)
	*/
}
