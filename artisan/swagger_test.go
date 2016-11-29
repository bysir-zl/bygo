package artisan

import (
	"testing"
	"github.com/deepzz0/go-com/log"
)

func TestSwagger(t *testing.T) {
	//walkApiFile("./swagger_te/","./swagger.json")
	x:=S("./","./swagger.json")
	log.Warn(x)
}
