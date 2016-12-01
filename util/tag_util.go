package util

import (
	"reflect"
	"sync"
)

type FieldTagMapper struct {
	// [tagName =>[FieldName=>tagValue]]
	mapData map[string]map[string]string
}

func (p *FieldTagMapper) GetFieldMapByTagName(tag string) (data map[string]string) {
	data = p.mapData[tag]
	return
}

var tagMapLock sync.RWMutex

// 从struct 取出 [tagName =>[fieldName=>tagValue]]
func newFieldTagMapper(i interface{}) (fieldTagMapper FieldTagMapper) {
	v := reflect.Indirect(reflect.ValueOf(i))

	fieldNum := v.NumField()

	reData := map[string]map[string]string{}

	for index := 0; index < fieldNum; index = index + 1 {
		f := v.Type().Field(index)
		x := EncodeTag(string(f.Tag))

		for tagKey, tagValue := range x {
			tagMapLock.RLock()
			if s := reData[tagKey]; s == nil {
				tagMapLock.RUnlock()
				tagMapLock.Lock()
				if s := reData[tagKey]; s == nil {
					reData[tagKey] = map[string]string{}
				}
				tagMapLock.Unlock()
			} else {
				tagMapLock.RUnlock()
			}
			reData[tagKey][f.Name] = tagValue
		}

	}

	fieldTagMapper = FieldTagMapper{}
	fieldTagMapper.mapData = reData

	return
}

var tagMapPoolLock sync.RWMutex
var mapperPool map[string]FieldTagMapper = map[string]FieldTagMapper{}

func GetTagMapperFromPool(i interface{}) FieldTagMapper {
	key := reflect.ValueOf(i).String()
	tagMapPoolLock.RLock()
	if s := mapperPool[key]; s.mapData == nil {
		tagMapPoolLock.RUnlock()
		tagMapPoolLock.Lock()
		if s := mapperPool[key]; s.mapData == nil {
			mapperPool[key] = newFieldTagMapper(i)
		}
		tagMapPoolLock.Unlock()
	} else {
		tagMapPoolLock.RUnlock()
	}

	return mapperPool[key]
}
