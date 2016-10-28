package util

import (
	"testing"
	"github.com/deepzz0/go-com/log"
)

func Test_Conv(b *testing.T) {
	sr := []byte("asg_as_bf")
	bs := SheXing2TuoFeng(sr)
	log.Print(string(bs))
}
