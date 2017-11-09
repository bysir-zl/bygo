package wxencrypter

import (
	"encoding/xml"
)

type EncryptedResponseXML struct {
	XMLName      xml.Name `xml:"xml"`
	TimeStamp    string
	Encrypt      string
	MsgSignature string
	Nonce        string
}

type EncryptedRequestXML struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string
	Encrypt    string
}

func ParseRequestXML(data []byte) (e EncryptedRequestXML, err error) {
	err = xml.Unmarshal(data, &e)
	if err != nil {
		err = ParseXmlError
	}
	return
}

func ParseResponseXML(data []byte) (e EncryptedResponseXML, err error) {
	err = xml.Unmarshal(data, &e)
	if err != nil {
		err = ParseXmlError
	}
	return
}

func GenerateResponseXML(encrypt, signature, timestamp, nonce string) (b []byte, err error) {
	e := EncryptedResponseXML{
		Nonce:        nonce,
		Encrypt:      encrypt,
		TimeStamp:    timestamp,
		MsgSignature: signature,
	}

	b, err = xml.Marshal(e)
	if err != nil {
		err = GenReturnXmlError
	}
	return
}
