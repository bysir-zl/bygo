package util

import (
	"reflect"
	"strconv"
	"strings"
	"fmt"
	"net/url"
	"github.com/deepzz0/go-com/log"
)

func EncodeTag(tag string) (data map[string]string) {
	data = map[string]string{}
	if tag == "" {
		return
	}
	for _, item := range strings.Split(tag, " ") {
		if item == "" {
			continue
		}
		key := strings.Split(item, ":")[0]
		value := strings.Split(item, "\"")[1]
		data[key] = value
	}

	return
}

func MapListToObjList(obj interface{}, mappers []map[string]interface{}, useTag string) {

	pointer := reflect.Indirect(reflect.ValueOf(obj))
	typer := pointer.Type().Elem()

	for _, mapper := range mappers {
		item := reflect.New(typer)
		MapToObj(item.Interface(), mapper, useTag)
		pointer.Set(reflect.Append(pointer, reflect.Indirect(item)))
	}
}

func ObjListToMapList(obj interface{}, useTag string) (mappers []map[string]interface{}) {
	mappers = []map[string]interface{}{}

	value := reflect.ValueOf(obj)
	for i := 0; i < value.Len(); i = i + 1 {
		item := value.Index(i)
		mappers = append(mappers, ObjToMap(item.Interface(), useTag))
	}
	return
}

// 根据map的key=>value设置Obj的field=>fieldValue
// 如果传了useTag,那么就会根据obj的Tag的useTag的值获取mapValue并填充到field上,
// 返回设置成功的Fields列表字段
func MapToObj(obj interface{}, mapper map[string]interface{}, useTag string) (fields []string) {
	pointer := reflect.Indirect(reflect.ValueOf(obj))
	typer := pointer.Type()
	fieldNum := pointer.NumField()

	var fieldNameToTagName map[string]string
	if useTag != "" {
		fieldTagMapper := GetTagMapperFromPool(obj)
		fieldNameToTagName = fieldTagMapper.GetFieldMapByTagName(useTag)
	}

	fields = []string{}
	for i := 0; i < fieldNum; i++ {
		field := pointer.Field(i)
		fieldName := typer.Field(i).Name
		key := fieldName

		if useTag != "" {
			// 根据指定的tag的key重新映射
			key = fieldNameToTagName[key]
			// 如果有逗号 比如 json:"password,omitempty" 则只取逗号前面的第一个
			key = strings.Split(key, ",")[0]
		}

		if value := mapper[key]; value != nil {
			if field.CanInterface() {
				setFieldValue(field, value)
				fields = append(fields, fieldName)
			}
		}
	}
	return
}

func MapStringToObj(obj interface{}, mapper map[string]string, useTag string) (fields []string) {
	mapper2 := map[string]interface{}{}
	for k, v := range mapper {
		mapper2[k] = v
	}
	return MapToObj(obj, mapper2, useTag)
}

func setFieldValue(field reflect.Value, value interface{}) {
	switch field.Interface().(type) {
	case bool:
		switch value.(type) {
		case bool:
			field.SetBool(value.(bool))
			break

		case string:
			s := value.(string)
			field.SetBool(s == "1" || strings.ToLower(s) == "true")
			break
		}
		break
	case string:
		strv := ""
		switch value.(type) {
		case string:
			strv = value.(string)
			break
		case []uint8:
			strv = string(value.([]uint8))
			break
		default:
			log.Print("not case type : " + field.Type().Name() + " is " + reflect.ValueOf(value).Type().Kind().String() + " in db , not " + field.Type().Kind().String())
			break
		}
		field.SetString(strv)
		break
	case int, int8, int16, int32, int64:
		var intv int64 = 0

		switch value.(type) {
		case int, int8, int32, int64:
			intv = intInterfaceToInt64(value)
		case float32:
			intv = int64(value.(float32))
		case float64:
			intv = int64(value.(float64))
		case string:
			intv, _ = strconv.ParseInt(value.(string), 10, 64)
		case []uint8:
			intv, _ = strconv.ParseInt(string(value.([]uint8)), 10, 64)
		}
		field.SetInt(intv)
		break
	case float32, float64:
		var flov float64 = 0
		switch value.(type) {
		case float32:
			flov = float64(value.(float32))
		case float64:
			flov = float64(value.(float64))
		case string:
			flov, _ = strconv.ParseFloat(value.(string), 64)
			break
		}
		field.SetFloat(flov)
		break
	default:
		println("not case type : " + field.Type().String())
		break
	}
}

func intInterfaceToInt64(value interface{}) int64 {
	switch value.(type) {
	case int:
		return int64(value.(int))
	case int8:
		return int64(value.(int8))
	case int32:
		return int64(value.(int32))
	case int64:
		return int64(value.(int64))
	}
	return 0
}

func ObjToMap(obj interface{}, useTag string) map[string]interface{} {
	pointer := reflect.Indirect(reflect.ValueOf(obj))
	typer := pointer.Type()

	fieldNum := pointer.NumField()

	var fieldNameToTagName map[string]string
	if useTag != "" {
		fieldTagMapper := GetTagMapperFromPool(obj)
		fieldNameToTagName = fieldTagMapper.GetFieldMapByTagName(useTag)
	}

	data := map[string]interface{}{}

	for i := 0; i < fieldNum; i++ {
		field := pointer.Field(i)
		key := typer.Field(i).Name

		if useTag != "" {
			// 根据指定的tag的key重新映射
			key = fieldNameToTagName[key]
			// 如果有逗号 比如 json:"password,omitempty" 则只取逗号前面的第一个
			key = strings.Split(key, ",")[0]
			// 有值才填充
			if key != "" {
				data[key] = field.Interface()
			}
		} else {
			data[key] = field.Interface()
		}

	}

	return data
}

//将map[string'key']string'value'  转换为map[value]key
func ReverseMap(ma map[string]string) (data map[string]string) {
	data = map[string]string{}

	for key, value := range ma {
		data[value] = key
	}

	return
}

//判断一个array每一个原始是不是都在map的value里
func ArrayInMapValue(min []string, m map[string]string) (has bool, msg string) {
	if min == nil || len(min) == 0 {
		return true, ""
	}
	lenMin := len(min)
	for minI := 0; minI < lenMin; minI = minI + 1 {
		has := false
		for _, value := range m {
			if value == min[minI] {
				has = true
			}
		}
		if !has {
			return false, min[minI]
		}
	}
	return true, ""
}

//获取map的keys
func GetMapKey(m map[string]string) (keys []string) {
	keys = []string{}

	for key, _ := range m {
		keys = append(keys, key)
	}

	return keys
}

//判断一个array每一个原始是不是都在map的key里
func ArrayInMapKey(min []string, m map[string]string) (has bool, msg string) {
	if min == nil || len(min) == 0 {
		has = true
		return
	}
	if m == nil {
		has = false
		return
	}
	lenMin := len(min)
	for minI := 0; minI < lenMin; minI = minI + 1 {
		has = false
		for key, _ := range m {
			if key == min[minI] {
				has = true
			}
		}
		if !has {
			msg = min[minI]
			return
		}
	}
	has = true
	return
}

func ArrayInArray(min []string, max []string) (has bool, msg string) {
	if min == nil || len(min) == 0 {
		return true, ""
	}

	lenMax := len(max)
	lenMin := len(min)
	for minI := 0; minI < lenMin; minI = minI + 1 {
		has := false
		for maxI := 0; maxI < lenMax; maxI = maxI + 1 {
			if max[maxI] == min[minI] {
				has = true
			}
		}
		if !has {
			return false, min[minI]
		}
	}
	return true, ""
}

// 判断item是否在数组里
// 如果数组为空则返回false
func ItemInArray(item string, max []string) (has bool) {

	if max == nil || len(max) == 0 {
		return false
	}

	lenMax := len(max)

	for maxI := 0; maxI < lenMax; maxI = maxI + 1 {
		if max[maxI] == item {
			return true
		}
	}

	return false
}

// 判断item是否在数组里
// 如果数组为空则返回false
func ItemInArrayInt(item int, max []int) (has bool) {

	if max == nil || len(max) == 0 {
		return false
	}

	lenMax := len(max)

	for maxI := 0; maxI < lenMax; maxI = maxI + 1 {
		if max[maxI] == item {
			return true
		}
	}

	return false
}

func IsEmptyValue(value interface{}) bool {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func EmptyObject(obj interface{}) {
	pointer := reflect.Indirect(reflect.ValueOf(obj))
	fieldNum := pointer.NumField()

	for i := 0; i < fieldNum; i++ {
		v := pointer.Field(i)
		switch v.Kind() {
		case reflect.String:
			v.SetString("")
			break
		case reflect.Bool:
			v.SetBool(false)
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v.SetInt(0)
			break
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v.SetUint(0)
			break
		case reflect.Float32, reflect.Float64:
			v.SetFloat(0)
			break
		case reflect.Interface, reflect.Ptr:
			break
		}
	}
}

func MapInterface2MapString(m map[string]interface{}) map[string]string {
	set := map[string]string{}

	for key, value := range m {
		switch value.(type) {
		case int64:
			i := value.(int64)
			set[key] = strconv.FormatInt(i, 10)
		case int32:
			i := int64(value.(int32))
			set[key] = strconv.FormatInt(i, 10)
		case int:
			i := int64(value.(int))
			set[key] = strconv.FormatInt(i, 10)
		case []byte:
			set[key] = string(value.([]byte))
		case string:
			set[key] = value.(string)
		case float64:
			set[key] = fmt.Sprintf("%f", value.(float64))
		default:
			log.Warn(key, " is not cased! :" + reflect.ValueOf(value).Type().String())
		}
	}
	return set

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
