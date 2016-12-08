package encoder

import (
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"crypto"
	"encoding/pem"
	"math/big"
	"errors"
	"io"
)

func RsaSign(origin, priKey []byte, hash crypto.Hash) ([]byte, error) {
	b, _ := pem.Decode(priKey)
	prikey, err := genPriKey(b.Bytes)
	if err != nil {
		return nil, err
	}

	h := hash.New()
	h.Write(origin)
	hashed := h.Sum(nil)

	sign, err := rsa.SignPKCS1v15(rand.Reader, prikey, hash, hashed)

	return sign, err
}

func RsaVerify(origin, sign, pubKey []byte, hash crypto.Hash) error {
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

func RsaPubDecrypt(origin, pubKey []byte) (data []byte, err error) {
	b, _ := pem.Decode(pubKey)
	pubkey, err := genPubKey(b.Bytes)
	if err != nil {
		return
	}

	data, err = pubKeyDecrypt(pubkey, origin)
	return
}

func RsaPriDecrypt(origin, priKey []byte) (data []byte, err error) {
	b, _ := pem.Decode(priKey)
	prikey, err := genPriKey(b.Bytes)
	if err != nil {
		return
	}

	data, err = rsa.DecryptPKCS1v15(rand.Reader, prikey, origin)
	return
}

func RsaPriEncrypt(origin, priKey []byte) (data []byte, err error) {
	b, _ := pem.Decode(priKey)
	prikey, err := genPriKey(b.Bytes)
	if err != nil {
		return
	}

	data, err = priKeyEncrypt(rand.Reader, prikey, origin)
	return
}

func RsaPubEncrypt(origin, pubKey []byte) (data []byte, err error) {
	b, _ := pem.Decode(pubKey)
	pubkey, err := genPubKey(b.Bytes)
	if err != nil {
		return
	}

	data, err = rsa.EncryptPKCS1v15(rand.Reader, pubkey, origin)
	return
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
	// 秘钥格式一般有两种 : PKCS#1,PKCS#8
	prkI, err := x509.ParsePKCS1PrivateKey(privateKey)
	if err != nil {
		prkI2, er := x509.ParsePKCS8PrivateKey(privateKey)
		if er != nil {
			return nil, er
		}
		prkI = prkI2.(*rsa.PrivateKey)
	}

	priKey = prkI
	return priKey, nil
}

// below code from github.com/qyxing/ostar
var (
	ErrDataToLarge = errors.New("message too long for RSA public key size")
	ErrDataLen = errors.New("data length error")
	ErrDataBroken = errors.New("data broken, first byte is not zero")
	ErrKeyPairDismatch = errors.New("data is not encrypted by the private key")
	ErrDecryption = errors.New("decryption error")
	ErrPublicKey = errors.New("get public key error")
	ErrPrivateKey = errors.New("get private key error")
)

func pubKeyDecrypt(pub *rsa.PublicKey, data []byte) ([]byte, error) {
	k := (pub.N.BitLen() + 7) / 8
	if k != len(data) {
		return nil, ErrDataLen
	}
	m := new(big.Int).SetBytes(data)
	if m.Cmp(pub.N) > 0 {
		return nil, ErrDataToLarge
	}
	m.Exp(m, big.NewInt(int64(pub.E)), pub.N)
	d := leftPad(m.Bytes(), k)
	if d[0] != 0 {
		return nil, ErrDataBroken
	}
	if d[1] != 0 && d[1] != 1 {
		return nil, ErrKeyPairDismatch
	}
	var i = 2
	for ; i < len(d); i++ {
		if d[i] == 0 {
			break
		}
	}
	i++
	if i == len(d) {
		return nil, nil
	}
	return d[i:], nil
}

/*私钥加密*/
func priKeyEncrypt(rand io.Reader, priv *rsa.PrivateKey, hashed []byte) ([]byte, error) {
	tLen := len(hashed)
	k := (priv.N.BitLen() + 7) / 8
	if k < tLen + 11 {
		return nil, ErrDataLen
	}
	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < k - tLen - 1; i++ {
		em[i] = 0xff
	}
	copy(em[k - tLen:k], hashed)
	m := new(big.Int).SetBytes(em)
	c, err := decrypt(rand, priv, m)
	if err != nil {
		return nil, err
	}
	copyWithLeftPad(em, c.Bytes())
	return em, nil
}

func encrypt(c *big.Int, pub *rsa.PublicKey, m *big.Int) *big.Int {
	e := big.NewInt(int64(pub.E))
	c.Exp(m, e, pub.N)
	return c
}

var bigZero = big.NewInt(0)
var bigOne = big.NewInt(1)

func decrypt(random io.Reader, priv *rsa.PrivateKey, c *big.Int) (m *big.Int, err error) {
	if c.Cmp(priv.N) > 0 {
		err = ErrDecryption
		return
	}
	var ir *big.Int
	if random != nil {
		var r *big.Int

		for {
			r, err = rand.Int(random, priv.N)
			if err != nil {
				return
			}
			if r.Cmp(bigZero) == 0 {
				r = bigOne
			}
			var ok bool
			ir, ok = modInverse(r, priv.N)
			if ok {
				break
			}
		}
		bigE := big.NewInt(int64(priv.E))
		rpowe := new(big.Int).Exp(r, bigE, priv.N)
		cCopy := new(big.Int).Set(c)
		cCopy.Mul(cCopy, rpowe)
		cCopy.Mod(cCopy, priv.N)
		c = cCopy
	}

	if priv.Precomputed.Dp == nil {
		m = new(big.Int).Exp(c, priv.D, priv.N)
	} else {
		m = new(big.Int).Exp(c, priv.Precomputed.Dp, priv.Primes[0])
		m2 := new(big.Int).Exp(c, priv.Precomputed.Dq, priv.Primes[1])
		m.Sub(m, m2)
		if m.Sign() < 0 {
			m.Add(m, priv.Primes[0])
		}
		m.Mul(m, priv.Precomputed.Qinv)
		m.Mod(m, priv.Primes[0])
		m.Mul(m, priv.Primes[1])
		m.Add(m, m2)

		for i, values := range priv.Precomputed.CRTValues {
			prime := priv.Primes[2 + i]
			m2.Exp(c, values.Exp, prime)
			m2.Sub(m2, m)
			m2.Mul(m2, values.Coeff)
			m2.Mod(m2, prime)
			if m2.Sign() < 0 {
				m2.Add(m2, prime)
			}
			m2.Mul(m2, values.R)
			m.Add(m, m2)
		}
	}
	if ir != nil {
		m.Mul(m, ir)
		m.Mod(m, priv.N)
	}

	return
}

func copyWithLeftPad(dest, src []byte) {
	numPaddingBytes := len(dest) - len(src)
	for i := 0; i < numPaddingBytes; i++ {
		dest[i] = 0
	}
	copy(dest[numPaddingBytes:], src)
}

func leftPad(input []byte, size int) (out []byte) {
	n := len(input)
	if n > size {
		n = size
	}
	out = make([]byte, size)
	copy(out[len(out) - n:], input)
	return
}

func modInverse(a, n *big.Int) (ia *big.Int, ok bool) {
	g := new(big.Int)
	x := new(big.Int)
	y := new(big.Int)
	g.GCD(x, y, a, n)
	if g.Cmp(bigOne) != 0 {
		return
	}
	if x.Cmp(bigOne) < 0 {
		x.Add(x, n)
	}
	return x, true
}