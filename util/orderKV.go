package util

import (
	"bytes"
	"net/url"
	"sort"
)

type OrderKV struct {
	keys   []string
	values [][]string
}

func (p *OrderKV) Add(key, value string) {
	if p.keys == nil {
		p.keys = []string{}
	}
	if p.values == nil {
		p.values = [][]string{}
	}
	if i := ArrayStringIndex(key, p.keys); i != -1 {
		if p.values[i] == nil {
			p.values[i] = []string{value}
		} else {
			p.values[i] = append(p.values[i], value)
		}
	} else {
		p.keys = append(p.keys, key)
		p.values = append(p.values, []string{value})
	}
}

func (p *OrderKV) Set(key, value string) {
	if p.keys == nil {
		p.keys = []string{}
	}
	if p.values == nil {
		p.values = [][]string{}
	}
	if i := ArrayStringIndex(key, p.keys); i != -1 {
		p.values[i] = []string{value}
	} else {
		p.keys = append(p.keys, key)
		p.values = append(p.values, []string{value})
	}
}

func ParseOrderKV(m map[string]string) OrderKV {
	o := OrderKV{}
	ks, vs := SortMap(m)
	for i, k := range ks {
		o.Set(k, vs[i])
	}
	return o
}

func (p *OrderKV) Map() map[string]string {
	set := map[string]string{}
	for i, k := range p.keys {
		set[k] = p.values[i][0]
	}
	return set
}

func (p OrderKV) String() string {
	return p.EncodeString()
}

func (p *OrderKV) MapMulti() map[string][]string {
	set := map[string][]string{}
	for i, k := range p.keys {
		set[k] = p.values[i]
	}
	return set
}

func (p *OrderKV) Keys() []string {
	return p.keys
}

func (p *OrderKV) Values() []string {
	set := make([]string, len(p.values))
	for i := range p.values {
		set[i] = p.values[i][0]
	}
	return set
}

func (p *OrderKV) Sort() {
	m := p.MapMulti()
	sort.Strings(p.keys)

	values := [][]string{}
	for _, k := range p.keys {
		values = append(values, m[k])
	}
	p.values = values
}

// See url.Values.Encode
func (p *OrderKV) Encode() []byte {
	if len(p.keys) == 0 {
		return []byte{}
	}
	var buf bytes.Buffer
	for i, k := range p.keys {
		k = url.QueryEscape(k)
		for _, v := range p.values[i] {
			v = url.QueryEscape(v)
			buf.WriteByte('&')
			buf.WriteString(k + "=" + v)
		}
	}
	return buf.Bytes()[1:]
}

func (p *OrderKV) EncodePhp() []byte {
	var buf bytes.Buffer
	for i, k := range p.keys {
		k = url.QueryEscape(k)
		multi := len(p.values[i]) != 1
		for _, v := range p.values[i] {
			v = url.QueryEscape(v)
			buf.WriteByte('&')
			if multi {
				buf.WriteString(k + "[]=" + v)
			} else {
				buf.WriteString(k + "=" + v)
			}
		}
	}
	return buf.Bytes()[1:]
}

func (p *OrderKV) EncodeString() string {
	return B2S(p.Encode())
}

func (p *OrderKV) EncodeStringWithoutEscape() string {
	var buf bytes.Buffer
	for i, k := range p.keys {
		for _, v := range p.values[i] {
			buf.WriteByte('&')
			buf.WriteString(k + "=" + v)
		}
	}
	return buf.String()[1:]
}

func (p *OrderKV) UrlValue() url.Values {
	u := url.Values(p.MapMulti())
	return u
}
