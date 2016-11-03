package util

import (
	"github.com/deepzz0/go-com/log"
	"testing"
)

func Test_Conv(b *testing.T) {
	sr := []byte("asg_as_bf")
	bs := SheXing2TuoFeng(sr)
	log.Print(string(bs))
}
