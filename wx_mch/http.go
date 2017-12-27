// 封装微信部分请求需要证书的HTTP请求

package wx_mch

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/bysir-zl/bygo/log"
	"io/ioutil"
	"net/http"
	"bytes"
)

var tlsConfig *tls.Config

func getTLSConfig() (*tls.Config, error) {
	if tlsConfig != nil {
		return tlsConfig, nil
	}

	// load cert
	cert, err := tls.LoadX509KeyPair(CertPemPath, KeyPemPath)
	if err != nil {
		log.Error("load wechat keys fail", err)
		return nil, err
	}

	// load root ca
	caData, err := ioutil.ReadFile(CAPemPath)
	if err != nil {
		log.Error("read wechat ca fail", err)
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
	return tlsConfig, nil
}

func SecurePost(url string, xmlContent []byte) ([]byte, error) {
	tlsConfig, err := getTLSConfig()
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: tr}
	rsp, err := client.Post(url, "text/xml", bytes.NewBuffer(xmlContent))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	rsp.Body.Close()

	return body, nil
}
