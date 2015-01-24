package gcm

import (
	"encoding/json"
	"encoding/xml"
	"log"
)

const (
	CCS_MESSAGE = `<message id="%s"><gcm xmlns="google:mobile:data">%s</gcm></message>`
)

type ccsMessage struct {
	XMLName xml.Name `xml:"message"`
	ID      string   `xml:"id,attr"`
	Body    string   `xml:"google:mobile:data gcm"`
}

type Message struct {
	To                       string      `json:"to"`
	MessageID                string      `json:"message_id"`
	Data                     interface{} `json:"data"`
	TTL                      int64       `json:"time_to_live"`
	DelayWhileIdle           bool        `json:"delay_while_idle"`
	DeliveryReceiptRequested bool        `json:"delivery_receipt_requested"`
}

func (this Message) Json() (string, error) {
	bytes, err := json.Marshal(this)
	if err != nil {
		log.Printf("Marshal ERROR: %+v ", err)
		return "", err
	}

	return string(bytes), err
}
