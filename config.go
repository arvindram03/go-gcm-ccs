package gcm

import "encoding/base64"

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
