package encoder

import (
	"testing"
	"github.com/bysir-zl/bygo/util"
	"github.com/deepzz0/go-com/log"
	"crypto"
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

	data:="appIdkqyzN5TtaUdUmutzualismIdKF_TEST01234mutualismUserId8e10bd9630c8efa62cc73acecb5e8da8time243"
	si,_:= RsaSign(util.S2B(data),util.S2B(pk),crypto.SHA1)
	log.Info(util.B2S(si))
}