package main

import (
	"log"
	"net"

	"crypto/tls"
)


const (
	GCM_API_KEY = "AIzaSyAzejjXJRtHsOJRXo1VXlrDELipdCqtfn0"
	GCM_CCS_ADDR = "gcm-preprod.googleapis.com:5236"
	TCP ="tcp"
)
func main() {
	tlsConfig := &tls.Config{
		ServerName: "gcm-preprod.googleapis.com",
	}
	tlsConn, err := GetTLSConn(GCM_CCS_ADDR, tlsConfig)
	defer tlsConn.Close()
	if err != nil {
		return
	}

	log.Printf("Conn: %+v", tlsConn.ConnectionState())
}

func GetTLSConn(addr string, tlsConfig *tls.Config) (tlsConn *tls.Conn,err error) {
	tcpConn, err := getTCPConn(addr)
	if err != nil {
		return nil, err
	}

	tlsConn = tls.Client(tcpConn, tlsConfig)
	err = tlsConn.Handshake()
	if err != nil {
		log.Printf("TLS handshake failed to ADDR:%s. ERROR: %+v",GCM_CCS_ADDR,  err)
		return nil, err
	}

	return tlsConn, nil
}

func getTCPConn(addr string) (conn net.Conn, err error) {
	conn, err = net.Dial(TCP, addr)
	if err != nil {
		log.Printf("Conenction failed to ADDR:%s. ERROR: %+v",addr,  err)
		return nil, err
	}
	return
}
