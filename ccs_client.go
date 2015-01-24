package gcm

import (
	"code.google.com/p/go-uuid/uuid"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net"
)

const (
	TCP = "tcp"
)

type GCMClient struct {
	CCSClient *CCSClient
}

func (this *GCMClient) NewClient(config Config) (client *CCSClient, err error) {
	tlsConfig := &tls.Config{
		ServerName: config.Host,
	}

	client = &CCSClient{}

	err = client.Init(tlsConfig, config)
	if err != nil {
		return nil, err
	}

	return
}

type CCSClient struct {
	tlsConn   *tls.Conn
	tlsConfig *tls.Config
	config    Config
	xmlStream *xml.Decoder
	jabberID  string
}

func (this *CCSClient) Init(tlsConfig *tls.Config, config Config) (err error) {
	this.tlsConfig = tlsConfig
	this.config = config
	err = this.initTLSConn()
	if err != nil {
		return
	}
	//log.Println("Initated TLS")
	err = this.tlsHandshake()
	if err != nil {
		return
	}
	//log.Println("Initated TLS Handshake")
	this.initXMLStream()
	//log.Println("Initated XML Stream")
	err = this.authenticate()
	if err != nil {
		return
	}
	//log.Println("Authenticated")
	return
}

func (this *CCSClient) Send(message Message) (err error) {
	jsonMessage, err := message.Json()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(this.tlsConn, CCS_MESSAGE, message.MessageID, jsonMessage)
	if err != nil {
		log.Panicln("ERROR sending message: %+v", err)
	}
	log.Println("Message Sent: %s ", fmt.Sprintf(CCS_MESSAGE, message.MessageID, jsonMessage))
	return
}

func (this *CCSClient) Recv() (err error) {
	xmlResponse, err := getXMLResponse(this.xmlStream)
	if err != nil {
		log.Println("ERROR Receiving: %+v", err)
	}
	ccsMessage := &ccsMessage{}
	if err = this.xmlStream.DecodeElement(ccsMessage, &xmlResponse); err != nil {
		return errors.New("ERROR UNMARSHALL <features>: " + err.Error())
	}
	log.Printf("Received Message: %s", ccsMessage.Body)
	return
}

func (this *CCSClient) Close() (err error) {
	if this.tlsConn != nil {
		err = this.tlsConn.Close()
	}
	return
}

func (this *CCSClient) initTLSConn() error {

	tcpConn, err := getTCPConn(this.config.FullAddress())
	if err != nil {
		return err
	}

	this.tlsConn = tls.Client(tcpConn, this.tlsConfig)
	return nil
}

func getTCPConn(addr string) (conn net.Conn, err error) {
	conn, err = net.Dial(TCP, addr)
	if err != nil {
		log.Printf("Conenction failed to ADDR:%s. ERROR: %+v", addr, err)
		return nil, err
	}
	return
}

func (this *CCSClient) tlsHandshake() error {
	err := this.tlsConn.Handshake()
	if err != nil {
		log.Printf("TLS handshake failed to ADDR:%s. ERROR: %+v", this.config.FullAddress(), err)
		return err
	}
	return nil
}

func (this *CCSClient) initXMLStream() {
	this.xmlStream = xml.NewDecoder(this.tlsConn)
}

func (this *CCSClient) authenticate() error {

	fmt.Fprintf(this.tlsConn, START_STREAM)
	xmlResponse, err := getXMLResponse(this.xmlStream)
	if err != nil {
		return err
	}

	if xmlResponse.Name.Space != XML_STREAM_NAMESPACE || xmlResponse.Name.Local != XML_STREAM_LOCAL_NAME {
		return fmt.Errorf("expected <stream> but got <%+v> in %+v", xmlResponse.Name.Local, xmlResponse.Name.Space)
	}

	f := new(streamFeatures)
	if err = this.xmlStream.DecodeElement(f, nil); err != nil {
		return errors.New("ERROR UNMARSHALL <features>: " + err.Error())
	}

	for _, mechanism := range f.Mechanisms.Mechanism {
		if mechanism == "PLAIN" {
			fmt.Fprintf(this.tlsConn, CLIENT_AUTH, this.config.GetEncodedKey())
			break
		}
	}

	xmlResponse, err = getXMLResponse(this.xmlStream)
	if err != nil {
		return err
	}

	if xmlResponse.Name.Space != XML_SASL_NAMESPACE || xmlResponse.Name.Local != XML_SASL_SUCCESS {
		return fmt.Errorf("expected <success> but got <%+v> in %+v", xmlResponse.Name.Local, xmlResponse.Name.Space)
	}

	response := &saslSuccess{}

	if err = this.xmlStream.DecodeElement(response, &xmlResponse); err != nil {
		return errors.New("ERROR UNMARSHALL <sasl success>: " + err.Error())
	}

	fmt.Fprintf(this.tlsConn, START_STREAM)
	xmlResponse, err = getXMLResponse(this.xmlStream)
	if err != nil {
		return err
	}

	f = new(streamFeatures)
	if err = this.xmlStream.DecodeElement(f, nil); err != nil {
		return errors.New("ERROR UNMARSHALL <features>: " + err.Error())
	}

	sessionID := uuid.New()

	fmt.Fprintf(this.tlsConn, IQ_BIND_REQUEST, sessionID)

	var iq clientIQ
	if err = this.xmlStream.DecodeElement(&iq, nil); err != nil {
		return errors.New("unmarshal <iq>: " + err.Error())
	}
	if &iq.Bind == nil {
		return errors.New("<iq> result missing <bind>")
	}

	this.jabberID = iq.Bind.Jid

	return nil
}
