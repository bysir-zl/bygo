package structs

import (
	"testing"
)

type X struct {
	B
	A A       `json:"a"`
	C []int     `json:"c"`
	F float64 `json:"f"`
}
type B struct {
	I int    `json:"i"`
	X string `json:"x"`
}
type A struct {
	A int `json:"a"`
}

func TestStruct2MapString(t *testing.T) {
	x := X{}
	m, err := Struct2MapString(x, "json")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", m)
}

func TestMap2Struct(t *testing.T) {
	x := X{}
	m := map[string]interface{}{
		"a": map[string]interface{}{"a": 1},
		"i": 1,
		"f": "123675567567.1232313777777",
		"x": 12312,
		"c": []int{1, 3, 4},
	}
	err := Map2Struct(m, &x, "json")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", x)
}
