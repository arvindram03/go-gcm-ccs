package gcm

import (
	"encoding/xml"
	"io"
)

const (
	START_STREAM = `<stream:stream to="gcm.googleapis.com" version="1.0" xmlns="jabber:client" 
        				xmlns:stream="http://etherx.jabber.org/streams">`
	CLIENT_AUTH = `<auth mechanism="PLAIN" xmlns="urn:ietf:params:xml:ns:xmpp-sasl">%s</auth>`
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
