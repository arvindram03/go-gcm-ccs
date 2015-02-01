package gcm

import (
	"encoding/xml"
)

const (
	CCS_MESSAGE = `<message id="%s"><gcm xmlns="google:mobile:data">%s</gcm></message>`
)

const (
	ACK      = "ack"
	NACK     = "nack"
	CONTROL  = "control"
	RECEIPT  = "receipt"
	UPSTREAM = ""
)

type ccsMessage struct {
	XMLName xml.Name `xml:"message"`
	ID      string   `xml:"id,attr"`
	Body    string   `xml:"google:mobile:data gcm"`
}

type CCSMessageResponse struct {
	MessageType      string      `json:"message_type"`
	ControlType      string      `json:"control_type"`
	MessageID        string      `json:"message_id"`
	From             string      `json:"from"`
	Data             interface{} `json:"data"`
	Category         string      `json:"category"`
	Error            string      `json:"error"`
	ErrorDescription string      `json:"error_description"`
}

type Message struct {
	To                       string      `json:"to"`
	MessageID                string      `json:"message_id"`
	Data                     interface{} `json:"data"`
	TTL                      int64       `json:"time_to_live"`
	DelayWhileIdle           bool        `json:"delay_while_idle"`
	DeliveryReceiptRequested bool        `json:"delivery_receipt_requested"`
}

type ACKMessage struct {
	To          string `json:"to"`
	MessageID   string `json:"message_id"`
	MessageType string `json:"message_type"`
}

type ReceiptMessage struct {
	MessageStatus        string `json:“message_status"`
	OriginalMessageID    string `json:“original_message_id”`
	DeviceRegistrationID string `json:“device_registration_id”`
}
