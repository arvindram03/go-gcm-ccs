package main

import (
	"log"
	"net"
	"crypto/tls"
)


const (
	GCM_API_KEY = "AIzaSyAzejjXJRtHsOJRXo1VXlrDELipdCqtfn0"
	GCM_CCS_HOST = "gcm-preprod.googleapis.com"
	GCM_CCS_PORT = "5236"
	GCM_CCS_ADDR = "gcm-preprod.googleapis.com:5236"
	TCP ="tcp"
)

type GCMCCSClient struct {
	TLSConn	*tls.Conn
	TLSConfig *tls.Config
}

func (this *GCMCCSClient) Init(host string) {
	this.TLSConfig = &tls.Config{
		ServerName: host,
	}
}

func (this *GCMCCSClient) GetTLSConn() (tlsConn *tls.Conn,err error) {
	tcpConn, err := getTCPConn(GCM_CCS_ADDR)
	if err != nil {
		return nil, err
	}

	this.TLSConn = tls.Client(tcpConn, this.TLSConfig)
	err = this.TLSConn.Handshake()
	if err != nil {
		log.Printf("TLS handshake failed to ADDR:%s. ERROR: %+v",GCM_CCS_ADDR,  err)
		return nil, err
	}

	return this.TLSConn, nil
}

func (this *GCMCCSClient) Close() (err error){
	if this.TLSConn != nil {
		err = this.TLSConn.Close()
	}
	return
}

func getTCPConn(addr string) (conn net.Conn, err error) {
	conn, err = net.Dial(TCP, addr)
	if err != nil {
		log.Printf("Conenction failed to ADDR:%s. ERROR: %+v",addr,  err)
		return nil, err
	}
	return
}

func main() {

	gcmCSSClient := &GCMCCSClient{}
	gcmCSSClient.Init(GCM_CCS_HOST)

	tlsConn, err := gcmCSSClient.GetTLSConn()
	defer tlsConn.Close()
	if err != nil {
		return
	}

	log.Printf("Conn: %+v", tlsConn.ConnectionState())
}
