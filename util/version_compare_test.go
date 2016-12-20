package util

import (
	"github.com/bysir-zl/bygo/log"
	"testing"
)

func TestVersionCompare(t *testing.T) {
	log.Debug("Test",VersionCompareBigger("4.1.6566","4.1.3"))
}

