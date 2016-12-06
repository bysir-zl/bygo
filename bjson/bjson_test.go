package bjson

import (
	"testing"
	"log"
)

func TestBjson_MapString(t *testing.T) {
	bs := []byte(`{"name":"zl","sex":1,"age":21,"data":{"hab":"code"}}`)
	bj, _ := New(bs)

	ms := bj.MapString()
	mi := bj.MapInterface()

	log.Print("ms: ", ms)
	log.Print("mi: ", mi)
	log.Printf("name: %s", bj.Pos("name").String())
	log.Printf("age: %d,%s", bj.Pos("age").Int(), bj.Pos("age").String())
	log.Printf("sex: %t,%d", bj.Pos("sex").Bool(), bj.Pos("sex").Int())
	log.Printf("hab: %s", bj.Pos("data").Pos("hab").String())

	log.Printf("E name: %d", bj.Pos("name").Int())

}
