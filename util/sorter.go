package util

import (
	"sort"
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

