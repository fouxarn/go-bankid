package bankid

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

// A RoundTripper executes a request passing the given req as the SOAP
// envelope body. The HTTP response is then de-serialized onto the resp
// object. Returns error in case an error occurs serializing req, making
// the HTTP request, or de-serializing the response.
type RoundTripper interface {
	RoundTrip(req, resp Message) error
	RoundTripSoap12(action string, req, resp Message) error
}

// Message is an opaque type used by the RoundTripper to carry XML
// documents for SOAP.
type Message interface{}

// Header is an opaque type used as the SOAP Header element in requests.
type Header interface{}

// AuthHeader is a Header to be encoded as the SOAP Header element in
// requests, to convey credentials for authentication.
type AuthHeader struct {
	Namespace string `xml:"xmlns:ns,attr"`
	Username  string `xml:"ns:username"`
	Password  string `xml:"ns:password"`
}

// Client is a SOAP client.
type Client struct {
	URL         string              // URL of the server
	Namespace   string              // SOAP Namespace
	Envelope    string              // Optional SOAP Envelope
	Header      Header              // Optional SOAP Header
	ContentType string              // Optional Content-Type (default text/xml)
	Config      *http.Client        // Optional HTTP client
	Pre         func(*http.Request) // Optional hook to modify outbound requests
}

// NewClient creates a new bankid-client with specified URL and certificates to access bankid-api
func NewClient(url string, cert tls.Certificate) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{cert},
		},
	}

	config := &http.Client{Transport: tr}

	return &Client{
		URL:       url,
		Namespace: "http://bankid.com/RpService/v4.0.0/types/",
		Config:    config,
	}
}

// NewTestClient creates a new bankid-client with test config and test certificates
func NewTestClient() (*Client, error) {
	cert, err := tls.X509KeyPair(testCert, testKey)
	if err != nil {
		return nil, err
	}

	return NewClient("https://appapi.test.bankid.com/rp/v4?wsdl", cert), nil
}

func doRoundTrip(c *Client, setHeaders func(*http.Request), in, out Message) error {
	req := &Envelope{
		EnvelopeAttr: c.Envelope,
		NSAttr:       c.Namespace,
		Header:       c.Header,
		Body:         Body{Message: in},
	}

	if req.EnvelopeAttr == "" {
		req.EnvelopeAttr = "http://schemas.xmlsoap.org/soap/envelope/"
	}
	if req.NSAttr == "" {
		req.NSAttr = c.URL
	}
	var b bytes.Buffer
	err := xml.NewEncoder(&b).Encode(req)
	if err != nil {
		return err
	}
	cli := c.Config
	if cli == nil {
		cli = http.DefaultClient
	}
	r, err := http.NewRequest("POST", c.URL, &b)
	if err != nil {
		return err
	}
	// log.Println("Sent: " + b.String())
	setHeaders(r)
	if c.Pre != nil {
		c.Pre(r)
	}
	resp, err := cli.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	/*
		if resp.StatusCode != http.StatusOK {
			// read only the first Mb of the body in error case
			limReader := io.LimitReader(resp.Body, 1024*1024)
			body, _ := ioutil.ReadAll(limReader)
			return fmt.Errorf("%q: %q", resp.Status, body)
		}*/
	body, _ := ioutil.ReadAll(resp.Body)
	// log.Println("Received: " + string(body))
	readerBody := bytes.NewReader(body)
	//return xml.NewDecoder(resp.Body).Decode(out)
	return xml.NewDecoder(readerBody).Decode(out)
}

// RoundTrip implements the RoundTripper interface.
func (c *Client) RoundTrip(in, out Message) error {
	headerFunc := func(r *http.Request) {
		ct := c.ContentType
		if ct == "" {
			ct = "text/xml"
		}
		r.Header.Set("Content-Type", ct)
		if in != nil {
			r.Header.Add("SOAPAction", fmt.Sprintf("%s/%s", c.Namespace, reflect.TypeOf(in).Elem().Name()))
		}
	}
	return doRoundTrip(c, headerFunc, in, out)
}

// RoundTripSoap12 implements the RoundTripper interface with SOAP1.2 Content-Type action extension
func (c *Client) RoundTripSoap12(action string, in, out Message) error {
	headerFunc := func(r *http.Request) {
		r.Header.Add("Content-Type", fmt.Sprintf("application/soap+xml; charset=utf-8; action=\"%s\"", action))
	}
	return doRoundTrip(c, headerFunc, in, out)
}

// Envelope is a SOAP envelope.
type Envelope struct {
	XMLName      xml.Name `xml:"soapenv:Envelope"`
	EnvelopeAttr string   `xml:"xmlns:soapenv,attr"`
	NSAttr       string   `xml:"xmlns:typ,attr"`
	Header       Message  `xml:"soapenv:Header"`
	Body         Body
}

// Body is the body of a SOAP envelope.
type Body struct {
	XMLName xml.Name `xml:"soapenv:Body"`
	Message Message
}
