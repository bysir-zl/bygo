package encoder

import (
	"crypto"
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/bygo/util"
	"testing"
)

func TestRsaSignWithSHA1(t *testing.T) {

	pk := `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCtNIJZoIcSkZ222evBDXc4AgOmGDjvDmat6W+0HnCXBGSH4TLW
Qj/MGFB4IbtZYOHHByvtwKffjlJaXjDF42ZYGFiv+125byzRIs86uxMITEOyGHy6
bJ8tx/O/XHTN8ZVUO2GxMvfIKCKYRJklUt0wBFKB1lOSZ8VhUdmTTwvfeQIDAQAB
AoGARFC2rRU00W0f0LQpWY6vHCcnO0bIhfmfJC8zgM6Ux+vSnwmC3KFqulxIuOlN
FNaylqbMe80GKZXgA4atJBAqGCTQgFMKjoXfMZBIUkntynHkg0sWgZNLB9vK+r89
QJwQbjMwiC7mLAzXLg8U2y4snYLDSFIhnIRmojzbJC8GWxECQQDcj4/dotJS6Yxs
185u94eXWJz3oNXARObJEIivm79bXONaaF+1kSInX99Z8O5RwWcEvb/kbMZe0yxI
vUvsfcOlAkEAyQkK2Prf+btByNWJc4NrFInlXe1Op4iQyibwfuACx0rywmJVSeDo
5vMLGD6Zeepg/Hu4ZMjpRTz5O3yXduRURQJAG9lYrgCIFAX/QCMDoslIapi6wR2i
v7MzfMHEsH+26r9QybKSGyfnKxeU6RNd1B7adiPLXflKFuENH2Yfdw3uLQJBAI/7
k/NXqvaXsUP//FPpOdYZ9VbSUdUXsGu4e+LC2fqWqUujVeZ12Rkf1UBmBVIWFaR/
j89PPhNC2lZKo8iZO+kCQQCK9ohwMsOEgjvA5uu23tI7ZBAZ8l7gYHHpIdAiDaCA
Tkwt7ECSYKOvXVPRv+TPy+My6HpSEsY4jE2HniYWv8z4
-----END RSA PRIVATE KEY-----
`

	//b,_:=pem.Decode(util.S2B(pk))
	//
	//_, err := x509.ParsePKCS1PrivateKey(b.Bytes)
	//log.Info(err)

	data := "appIdkqyzN5TtaUdUmutzualismIdKF_TEST01234mutualismUserId8e10bd9630c8efa62cc73acecb5e8da8time243"
	si, _ := RsaSign(util.S2B(data), util.S2B(pk), crypto.SHA1)
	log.Info("Test",util.B2S(si))
}

var pubkey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDC0zvqmjCPYqS/AeAFz376AUmh
3WLWt5JksJDV5o+m9DOZeQKTt94AteO2IvfurUfg/yA0V3Q5blu/Fdt4DWgzcfPH
PbKIlMA8tXCUK3KVKQe+FaNNyrgdVAL2TQY6njv4aT9k/S4ABjOnQEbUwsPlLNVS
8Q6MevwYKyOFTGT8dQIDAQAB
-----END PUBLIC KEY-----
`
var prikey = `-----BEGIN PRIVATE KEY-----
MIICeAIBADANBgkqhkiG9w0BAQEFAASCAmIwggJeAgEAAoGBAMLTO+qaMI9ipL8B
4AXPfvoBSaHdYta3kmSwkNXmj6b0M5l5ApO33gC147Yi9+6tR+D/IDRXdDluW78V
23gNaDNx88c9soiUwDy1cJQrcpUpB74Vo03KuB1UAvZNBjqeO/hpP2T9LgAGM6dA
RtTCw+Us1VLxDox6/BgrI4VMZPx1AgMBAAECgYEAu26RNDDHCwsxx/k71xs647ad
ajYcwsm0813S2ZaJGWvSwJHk4sx/rltPCYk20c6vWkzYZMLGNAJyDbIvhJ4RYcvk
wFAYIvPUJzFTsLTZBGoXCdifGdyR07uKL7U+HHKUSJ3+DjgTKrzxO0qao88KDxwz
IYb0lzO8ms87KBvCnnkCQQD9K2ihpTZd4CTu/r1WV4QIJQES5NjwLfRh6uRgc23D
+ovtY2e+Mwcb3g0auFssM5uro/RSF+13V+0LVYX8SEa3AkEAxQDWxh1WtplNTtFD
dW4/o28HbcCpzQlXtNEXaGVqb9t3LW2t0YjDAsLWs8WWXl66FCIOLbSqXYtoHkc1
COJKMwJBANDInrZICIjsk6jhPgXZkJIi6jrJrbqNO3ARBZwhNVGc6w6vntu1O1SZ
EBeMF+xg9y1avd+Byh1UzrE9K4z9kgsCQGlKaH/cYGMZjlMIz0gtE4AzMEI9jcNT
MfgnJJ6cTYXZQ1oZW6Q4txl7ryrH+PUZJdTq2q8c900l3BEKt9K2tzcCQQC9AMip
JiC46gjg847woLLlp8XOoAfFyGsUSqe3jmq2bQ/i7AyP3msoV5kzIYzrHNSdt4j6
DYmoJ+jaa5wZ7qNG
-----END PRIVATE KEY-----
`

func Test_RsapubDec(t *testing.T) {
	data := []byte("123456asd")

	enc, err := RsaPriEncrypt(data, []byte(prikey))
	if err != nil {
		log.Warn("Test",err)
	}

	dec, err := RsaPubDecrypt(enc, []byte(pubkey))
	if err != nil {
		log.Warn("Test",err)
	}

	log.Info("Test",string(dec)) // 123456asd
}