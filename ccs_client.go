package gcm

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	TCP = "tcp"
)

type Config struct {
	Host string
	Port string
	//GCM API Key
	Username string
	//GCM Project Number
	Password string
}

func (this Config) FullAddress() string {
	return this.Host + ":" + this.Port
}

func (this Config) GetEncodedKey() string {
	hexString := "\x00" + this.Username + "@" + this.Host + "\x00" + this.Password
	return base64.StdEncoding.EncodeToString([]byte(hexString))
}

type CCSClient struct {
	tlsConn   *tls.Conn
	tlsConfig *tls.Config
	config    Config
	xmlStream *xml.Decoder
}

func (this *CCSClient) Init(tlsConfig *tls.Config, config Config) (err error) {
	this.tlsConfig = tlsConfig
	this.config = config
	err = this.initTLSConn()
	if err != nil {
		return
	}

	err = this.tlsHandshake()
	if err != nil {
		return
	}

	this.initXMLStream()

	err = this.authenticate()
	if err != nil {
		return
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

	fmt.Printf("%+v", response)

	return nil
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

func (this *CCSClient) Close() (err error) {
	if this.tlsConn != nil {
		err = this.tlsConn.Close()
	}
	return
}

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
