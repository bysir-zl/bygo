package util

import (
	"log"
	"testing"
)

func TestMapToObj(t *testing.T) {
	m := map[string]interface{}{
		"name": "bysir",
		"sex":  true,
		"age":  21,
	}

	type S struct {
		Name string `json:"name"`
		Sex  int    `json:"sex"`
		Age  int    `json:"age"`
	}
	s := S{}
	b:= struct {
		S
		Name string `json:"name"`
	}{}

	MapToObj(&b, m, "json")
	log.Printf("%+v", s)
}

func TestMapToObj2(t *testing.T) {
	m := map[string]interface{}{
		"name": "bysir",
		"sex":  1,
		"age":  21,
	}

	type INT int

	s := struct {
		Name string `json:"name"`
		Sex  INT    `json:"sex"`
		Age  int    `json:"age"`
	}{}

	MapToObj(&s, m, "json")
	log.Printf("%+v", s)
}

func TestMapList(t *testing.T) {
	m := []map[string]interface{}{{
		"name": "bysir",
		"sex":  true,
		"age":  21,
	}}
	var s []*struct {
		Name string `json:"name"`
		Sex  int    `json:"sex"`
		Age  int    `json:"age"`
	}
	MapListToObjList(&s, m, "json")
	log.Printf("%+v", s)
}

func BenchmarkMapToObj(b *testing.B) {
	m := map[string]interface{}{
		"Name": "bysir",
		"Sex":  true,
		"Age":  21,
		"Baba": 21,
		"Mama": 21,
	}

	s := struct {
		Name string `json:"Name"`
		Sex  int    `json:"Sex"`
		Age  int    `json:"Age"`
		Baba int    `json:"Baba"`
		Mama int    `json:"Mama"`
	}{}

	for i := 0; i < b.N; i++ {
		// 4603 ns/op
		MapToObj(&s, m, "json")
		// 2221 ns/op
		//MapToObj(&s, m, "")
	}

	log.Printf("%+v", s)
}
