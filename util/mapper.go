package util

import (
	"sort"
	"net/url"
)

func SortMap(m map[string]string) (key []string, value []string) {
	key = GetMapKey(m)
	sort.Strings(key)
	value = []string{}
	for _, k := range key {
		value = append(value, m[k])
	}
	return
}

func FilterEmpty(m map[string]string) {
	for k, v := range m {
		if v == "" {
			delete(m, k)
		}
	}
}

func Map2UrlValues(m map[string]string) url.Values {
	v := url.Values{}
	for key, value := range m {
		v.Add(key, value)
	}
	return v
}

func CopyMapString(m map[string]string) map[string]string {
	set := map[string]string{}
	for key, value := range m {
		set[key] = value
	}
	return set
}

func FilterMapString(m map[string]string, keys ...string) {
	for k := range m {
		if !ItemInArray(k, keys) {
			delete(m, k)
		}
	}
}

func FilterMapByFun(m map[string]string, fun func(s string) string, keys ...string) {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			m[k] = fun(v)
		}
	}
}


//将map[string'key']string'value'  转换为map[value]key
func ReverseMap(ma map[string]string) (data map[string]string) {
	if ma == nil {
		return
	}
	data = map[string]string{}

	for key, value := range ma {
		data[value] = key
	}

	return
}