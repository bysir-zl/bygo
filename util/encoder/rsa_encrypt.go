package encoder

import (
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"crypto"
	"encoding/pem"
)

func RsaSign(origin, priKey []byte,hash crypto.Hash) ([]byte, error) {
	b, _ := pem.Decode(priKey)
	prikey, err := genPriKey(b.Bytes)
	if err != nil {
		return nil, err
	}

	h := hash.New()
	h.Write(origin)
	hashed := h.Sum(nil)

	sign, err := rsa.SignPKCS1v15(rand.Reader, prikey,hash, hashed)

	return sign, err
}

func RsaVerify(origin, sign, pubKey []byte,hash crypto.Hash) error {
	b, _ := pem.Decode(pubKey)
	pubkey, err := genPubKey(b.Bytes)
	if err != nil {
		return err
	}

	h := hash.New()
	h.Write(origin)
	hashed := h.Sum(nil)

	err = rsa.VerifyPKCS1v15(pubkey, hash, hashed, sign)
	return err
}

func genPubKey(publicKey []byte) (*rsa.PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

func genPriKey(privateKey []byte) (*rsa.PrivateKey, error) {
	var priKey *rsa.PrivateKey
	var err error
	//prkI, err := x509.ParsePKCS8PrivateKey(privateKey)
	prkI, err := x509.ParsePKCS1PrivateKey(privateKey)
	//prkI, err := x509.ParsePKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	//priKey = prkI.(*rsa.PrivateKey)
	priKey = prkI
	return priKey, nil
}

