package util

import (
	"reflect"
	"strconv"
	"strings"
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
	if mapper == nil || len(mapper) == 0 {
		return
	}
	pointer := reflect.Indirect(reflect.ValueOf(obj))
	var tag2field = map[string]string{}
	if useTag != "" {
		fieldTagMapper := GetTagMapperFromPool(obj)
		for k, v := range fieldTagMapper.GetFieldMapByTagName(useTag) {
			v = strings.Split(v, ",")[0]
			tag2field[v] = k
		}
	}

	fields = []string{}
	for fieldName, value := range mapper {
		if useTag != "" {
			fieldName = tag2field[fieldName]
		}
		field := pointer.FieldByName(fieldName)
		if field.IsValid() && field.CanInterface() {
			setFieldValue(field, value)
			fields = append(fields, fieldName)
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
		b, ok := Interface2Bool(value, false)
		if ok {
			field.SetBool(b)
		}
	case string:
		s, ok := Interface2String(value, false)
		if ok {
			field.SetString(s)
		}
	case int, int8, int16, int32, int64:
		i, ok := Interface2Int(value, false)
		if ok {
			field.SetInt(i)
		}
	case float32, float64:
		f, ok := Interface2Float(value, false)
		if ok {
			field.SetFloat(f)
		}
	default:
		println("not case type : " + field.Type().String())
		break
	}
}

func Interface2Int(value interface{}, strict bool) (v int64, ok bool) {
	switch value.(type) {
	case int:
		v, ok = int64(value.(int)), true
	case int8:
		v, ok = int64(value.(int8)), true
	case int32:
		v, ok = int64(value.(int32)), true
	case int64:
		v, ok = int64(value.(int64)), true
	}
	if ok {
		return
	}
	if strict {
		if !ok {
			return
		}
	}

	switch value.(type) {
	case string, []uint8:
		s, _ := Interface2String(value, true)
		i, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			v, ok = i, true
		}
	case float32, float64:
		f, _ := Interface2Float(value, true)
		v, ok = int64(f), true
	case bool:
		if value.(bool) {
			v = 1
		} else {
			v = 0
		}
		ok = true
	}
	return
}

func Interface2Bool(value interface{}, strict bool) (v bool, ok bool) {
	if strict {
		v, ok = value.(bool)
		return
	}
	switch value.(type) {
	case bool:
		v, ok = value.(bool), true
	case int8, int, int32, int64:
		i, _ := Interface2Int(value, true)
		v, ok = i == 1, true
	case float32, float64:
		i, _ := Interface2Float(value, true)
		v, ok = i == 1, true
	case string, []uint8:
		s, _ := Interface2String(value, true)
		s = strings.ToLower(s)
		v, ok = s == "1" || s == "true", true
	}

	return
}

func Interface2Float(value interface{}, strict bool) (v float64, ok bool) {
	switch value.(type) {
	case float32:
		v, ok = float64(value.(float32)), true
	case float64:
		v, ok = float64(value.(float64)), true
	}
	if ok {
		return
	}
	if strict {
		if !ok {
			return
		}
	}

	switch value.(type) {
	case int, int8, int32, int64:
		i, ok := Interface2Int(value, true)
		if ok {
			v, ok = float64(i), true
		}
	case string:
		f, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			return 0, false
		}
		return f, true
	case []uint8:
		f, err := strconv.ParseFloat(string(value.([]uint8)), 64)
		if err != nil {
			return 0, false
		}
		return f, true
	}
	return
}

func Interface2String(value interface{}, strict bool) (v string, ok bool) {
	switch value.(type) {
	case string:
		v, ok = value.(string), true
	case []uint8:
		v, ok = string(value.([]uint8)), true
	}

	if ok {
		return
	}
	if strict {
		if !ok {
			return
		}
	}

	switch value.(type) {
	case int64, int8, int32, int:
		i, _ := Interface2Int(value, true)
		v, ok = strconv.FormatInt(i, 10), true
	case float64, float32:
		f, _ := Interface2Float(value, true)
		v, ok = strconv.FormatFloat(f, 'f', -1, 64), true
	case bool:
		v, ok = strconv.FormatBool(value.(bool)), true
	}
	return
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


//判断一个array每一个原始是不是都在map的value里
func ArrayInMapValue(min []string, m map[string]string) (has bool, msg string) {
	if min == nil || len(min) == 0 {
		return true, ""
	}
	lenMin := len(min)
	for minI := 0; minI < lenMin; minI = minI + 1 {
		_, has = m[min[minI]]
		if !has {
			return false, min[minI]
		}
	}
	return true, ""
}

//获取map的keys
func GetMapKey(m map[string]string) (keys []string) {
	keys = []string{}
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

//判断一个array每一个元素是不是都在map的key里
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
	for minI := 0; minI < lenMin; minI++ {
		_, has = m[min[minI]]
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
	return ArrayStringIndex(item, max) != -1
}

func ArrayStringIndex(item string, max []string) (index int) {
	index = -1
	if max == nil || len(max) == 0 {
		return
	}
	for i, l := 0, len(max); i < l; i++ {
		if max[i] == item {
			index = i
			return
		}
	}
	return
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
		v, ok := Interface2String(value, false)
		if ok {
			set[key] = v
		}
	}
	return set
}