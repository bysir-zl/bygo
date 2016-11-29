package uuid

import (
	"testing"
	"github.com/deepzz0/go-com/log"
	"math/rand"
)

func TestRand( t *testing.T) {
	log.Print(Rand().Hex())
	log.Print(Rand().Hex())
	log.Print(Rand().Hex())
	log.Print(Rand().Hex())
	log.Print(Rand().Hex())
	log.Print(Rand().Hex())
	log.Print(Rand().Hex())
	log.Print(Rand().Hex())
	log.Print(Rand().Hex())
	log.Print(Rand().Hex())
	log.Print(Rand().Hex())
	log.Print(Rand().Hex())


	log.Print(rand.Int31n(256))
	log.Print(rand.Int31n(256))
	log.Print(rand.Int31n(256))
	log.Print(rand.Int31n(256))
	log.Print(rand.Int31n(256))
	log.Print(rand.Int31n(256))
	log.Print(rand.Int31n(256))
	log.Print(rand.Int31n(256))
	log.Print(rand.Int31n(256))
	log.Print(rand.Int31n(256))
	log.Print(rand.Int31())
	log.Print(rand.Int31())
	log.Print(rand.Int31())
	log.Print(rand.Int31())
	log.Print(rand.Int31())
	log.Print(rand.Int31())
	log.Print(rand.Int31())
	log.Print(rand.Int31())
	log.Print(rand.Int31())
	log.Print(rand.Int31())
	log.Print(rand.Int31())
	log.Print(rand.Int31())
	log.Print(rand.Int31())
	log.Print(rand.Int31n(256))
	log.Print(rand.Int31n(256))
}
