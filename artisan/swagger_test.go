package artisan

import (
	"github.com/bysir-zl/bygo/log"
	"testing"
)

func TestSwagger(t *testing.T) {
	//walkApiFile("./swagger_te/","./swagger.json")
	x := S("./", "./swagger.json")
	log.Verbose("Test", x)
}
