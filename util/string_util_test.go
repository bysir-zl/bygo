package util

import (
	"github.com/bysir-zl/bygo/log"
	"testing"
)

func Test_Conv(b *testing.T) {
	sr := []byte("asg_as_bf")
	bs := SheXing2TuoFeng(sr)
	log.Info("Test",string(bs))
	bs = TuoFeng2SheXing(bs)
	log.Info("Test",string(bs))
}
