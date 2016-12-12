package encoder

import (
	"testing"
	"golang.org/x/crypto/scrypt"
	"lib.com/deepzz0/go-com/log"
)

func TestSha256(t *testing.T) {

	pwd:=[]byte("123123")

	sb,err:=scrypt.Key(pwd,[]byte("salt01351565644564"), 16384, 8, 1, 32)
	if err != nil {
		log.Warn(err)
	}
	log.Print(Base64EncodeString(string(sb)))
}