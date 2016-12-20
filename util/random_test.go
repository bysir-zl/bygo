package util

import (
	"log"
	"testing"
)

func TestRandom(t *testing.T) {
	log.Print(Rand(0,5))
	log.Print(Rand(0,5))
	log.Print(Rand(0,5))
	log.Print(Rand(0,5))
	log.Print(Rand(0,5))
	log.Print(Rand(0,5))
	log.Print(Rand(1,5))


}

