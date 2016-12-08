package util

import (
	"testing"
	"kuaifa.com/kuaifa/kuaifa-api-service/tool/log"
)

func TestOrderKV(t *testing.T) {
	o := OrderKV{}
	o.Add("a", "1")
	o.Add("a", "2")
	o.Set("c", "1")
	o.Add("b", "1")
	o.Set("b", "2")
	o.Add("c", "2")

	o.Sort()

	log.Print(string(o.Encode())) // a=1&a=2&b=2&c=1&c=2
}
