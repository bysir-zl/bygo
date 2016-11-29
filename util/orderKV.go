package util

import (
	"bytes"
	"net/url"
)

type OrderKV struct {
	keys   []string
	values []string
}

func (p *OrderKV) Add(key, value string) {
	if p.keys == nil {
		p.keys = []string{}
	}
	if p.values == nil {
		p.values = []string{}
	}
	p.keys = append(p.keys, key)
	p.values = append(p.values, value)
}
func (p *OrderKV) Map() map[string]string {
	set := map[string]string{}
	for i, k := range p.keys {
		set[k] = p.values[i]
	}
	return set
}
func (p *OrderKV) Keys() []string {
	return p.keys
}
func (p *OrderKV) Values() []string {
	return p.values
}
func (p *OrderKV) QueryString() string {
	var buf bytes.Buffer
	for i, k := range p.keys {
		buf.WriteByte('&')
		k = url.QueryEscape(k)
		v := url.QueryEscape(p.values[i])
		buf.WriteString(k + "=" + v)
	}
	return buf.String()[1:]
}
func (p *OrderKV) UrlValue() url.Values {
	set := url.Values{}
	for i, k := range p.keys {
		set.Add(k, p.values[i])
	}
	return set
}