package gcm

import (
	"encoding/xml"
	"io"
)

const (
	START_STREAM = `<stream:stream to="gcm.googleapis.com" version="1.0" xmlns="jabber:client" 
        				xmlns:stream="http://etherx.jabber.org/streams">`
	CLIENT_AUTH = `<auth mechanism="PLAIN" xmlns="urn:ietf:params:xml:ns:xmpp-sasl">%s</auth>`

	IQ_BIND_REQUEST = `<iq type="set" id="%s"><bind xmlns="urn:ietf:params:xml:ns:xmpp-bind"></bind></iq>\n`
)

const (
	XML_STREAM_NAMESPACE  = "http://etherx.jabber.org/streams"
	XML_STREAM_LOCAL_NAME = "stream"

	XML_SASL_NAMESPACE = "urn:ietf:params:xml:ns:xmpp-sasl"
	XML_SASL_SUCCESS   = "success"
)

type streamFeatures struct {
	XMLName    xml.Name `xml:"http://etherx.jabber.org/streams features"`
	Mechanisms saslMechanisms
}

type saslMechanisms struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	Mechanism []string `xml:"mechanism"`
}

type saslSuccess struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
}

type clientIQ struct { // info/query
	XMLName xml.Name `xml:"jabber:client iq"`
	From    string   `xml:",attr"`
	ID      string   `xml:",attr"`
	To      string   `xml:",attr"`
	Type    string   `xml:",attr"` // error, get, result, set
	Error   clientError
	Bind    bindBind
}

type clientError struct {
	XMLName xml.Name `xml:"jabber:client error"`
	Code    string   `xml:",attr"`
	Type    string   `xml:",attr"`
	Any     xml.Name
	Text    string
}

type bindBind struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	Resource string
	Jid      string `xml:"jid"`
}

func getXMLResponse(xmlStream *xml.Decoder) (xml.StartElement, error) {
	for {
		token, err := xmlStream.Token()
		if err != nil && err != io.EOF {
			return xml.StartElement{}, err
		}
		switch tokenType := token.(type) {
		case xml.StartElement:
			return tokenType, nil
		}
	}
}
