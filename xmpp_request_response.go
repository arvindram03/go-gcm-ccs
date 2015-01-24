package gcm

import (
	"encoding/xml"
)

const (
	START_STREAM = `<stream:stream to="gcm.googleapis.com" version="1.0" xmlns="jabber:client" 
        				xmlns:stream="http://etherx.jabber.org/streams">`
	//ADQ5MTQ1NTI0Njk5OEBnY20uZ29vZ2xlYXBpcy5jb20AQUl6YVN5QXplampYSlJ0SHNPSlJYbzFWWGxyREVMaXBkQ3F0Zm4w
	CLIENT_AUTH = `<auth mechanism="PLAIN" xmlns="urn:ietf:params:xml:ns:xmpp-sasl">%s</auth>`
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
