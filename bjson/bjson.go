package bjson

import (
	"encoding/json"
	"strconv"
	"reflect"
	"log"
	"strings"
)

type Bjson struct {
	self interface{}
}

func New(data []byte) (*Bjson, error) {
	b := Bjson{
		self:0,
	}
	err := json.Unmarshal(data, &b.self)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (p *Bjson) MapString() map[string]string {
	if p.self == nil {
		return nil
	}
	if i, ok := p.self.(map[string]interface{}); ok {
		return mapInterface2MapString(i)
	}
	return nil
}

func (p *Bjson) MapInterface() map[string]interface{} {
	if p.self == nil {
		return nil
	}
	if i, ok := p.self.(map[string]interface{}); ok {
		return i
	}
	return nil
}

func (p *Bjson) String() string {
	if p.self == nil {
		return ""
	}
	s, _ := interface2String(p.self)
	return s
}

func (p *Bjson) Bool() bool {
	if p.self == nil {
		return false
	}
	s, _ := interface2String(p.self)
	s = strings.ToLower(s)
	return s == "true" || s == "1"
}

func (p *Bjson) Int() int {
	if p.self == nil {
		return 0
	}
	if i, ok := p.self.(int); ok {
		return i
	}
	s, _ := interface2String(p.self)
	s = strings.Split(s, ".")[0]
	i, _ := strconv.Atoi(s)
	return i
}

func (p *Bjson) Float() float64 {
	if p.self == nil {
		return 0
	}
	if i, ok := p.self.(float64); ok {
		return i
	}
	if i, ok := p.self.(float32); ok {
		return float64(i)
	}
	s, _ := interface2String(p.self)
	i, _ := strconv.ParseFloat(s, 64)
	return i
}

func (p *Bjson)Pos(key string) *Bjson {
	if p.self == nil {
		return p
	}

	b := Bjson{}
	if i, ok := p.self.(map[string]interface{}); ok {
		b.self = i[key]
	}
	return &b
}

func interface2String(value interface{}) (string, bool) {
	switch value.(type) {
	case int64:
		i := value.(int64)
		return strconv.FormatInt(i, 10), true
	case int32:
		i := int64(value.(int32))
		return strconv.FormatInt(i, 10), true
	case int:
		i := int64(value.(int))
		return strconv.FormatInt(i, 10), true
	case []byte:
		return string(value.([]byte)), true
	case string:
		return value.(string), true
	case float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64), true
	case float32:
		return strconv.FormatFloat(float64(value.(float32)), 'f', -1, 64), true
	case bool:
		return strconv.FormatBool(value.(bool)), true
	case map[string]interface{}:
		return "", true
	}
	return "", false
}

func mapInterface2MapString(m map[string]interface{}) map[string]string {
	set := map[string]string{}

	for key, value := range m {
		v, ok := interface2String(value)
		if !ok {
			// debug
			log.Print(key + " is not cased!: " + reflect.TypeOf(value).String())
		} else {
			set[key] = v
		}
	}
	return set
}