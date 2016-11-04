package util

import (
	"testing"
	"lib.com/deepzz0/go-com/log"
)

func TestRandom(t *testing.T) {
	log.Print(Rand(0,5))
}

