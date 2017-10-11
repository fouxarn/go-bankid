package bankid

import (
	"encoding/xml"
)

type soapEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Header  *soapHeader
	Body    soapBody
}

type soapHeader struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`

	Items []interface{} `xml:",omitempty"`
}

type soapBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Fault   *soapFault  `xml:",omitempty"`
	Content interface{} `xml:",omitempty"`
}

type soapFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`

	Code   string     `xml:"faultcode,omitempty"`
	String string     `xml:"faultstring,omitempty"`
	Actor  string     `xml:"faultactor,omitempty"`
	Detail soapDetail `xml:",omitempty"`
}

type soapDetail struct {
	XMLName xml.Name `xml:"detail"`

	FaultStatus string `xml:"faultStatus"`
	Description string `xml:"detailedDescription"`
}

func (f *soapFault) Error() string {
	return f.String
}
