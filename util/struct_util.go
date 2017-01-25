package util

import (
	"fmt"
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

func MapListToObjList(obj interface{}, mappers []map[string]interface{}, useTag string) (errInfo string) {
	objValue := indirect(reflect.ValueOf(obj), false)
	item := GetElemInterface(reflect.ValueOf(obj))
	var e string
	for _, mapper := range mappers {
		iv:=reflect.New(reflect.TypeOf(item))
		_, e = MapToObj(iv.Interface(), mapper, useTag)
		objValue.Set(reflect.Append(objValue, iv.Elem()))
	}
	return e
}

func MapToObj(obj interface{}, mapper map[string]interface{}, useTag string) (fields []string, errInfo string) {
	if mapper == nil || len(mapper) == 0 {
		return
	}
	//log.Info("x2", reflect.TypeOf(obj))
	objValue := indirect(reflect.ValueOf(obj), false)
	//log.Info("x", objValue.Type())
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
		field := objValue.FieldByName(fieldName)
		if field.IsValid() && field.CanInterface() && field.CanSet() {
			err := setValue(field, value)
			if err != nil {
				errInfo = "field(" + fieldName + ") " + err.Error()
			} else {
				fields = append(fields, fieldName)
			}
		}
	}

	return
}

// copy from decode.go, i can't understand ...
//
// indirect walks down v allocating pointers as needed,
// until it gets to a non-pointer.
// if it encounters an Unmarshaler, indirect stops and returns that.
// if decodingNull is true, indirect stops at the last pointer so it can be set to nil.
func indirect(v reflect.Value, decodingNull bool) (reflect.Value) {
	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.Elem().Kind() != reflect.Ptr && decodingNull && v.CanSet() {
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 {

		}
		v = v.Elem()
	}
	return v
}

func GetElemInterface(v reflect.Value) interface{} {
	xx := indirect(v, false).Type()

	if xx.Kind() == reflect.Ptr {
		xx = xx.Elem()
	}
	if xx.Kind() == reflect.Slice {
		xx = xx.Elem()
	}

	return reflect.New(xx).Elem().Interface()
}

// 根据map的key=>value设置Obj的field=>fieldValue
// 如果传了useTag,那么就会根据obj的Tag的useTag的值获取mapValue并填充到field上,
// 返回设置成功的Fields列表字段
func ObjListToMapList(obj interface{}, useTag string) (mappers []map[string]interface{}) {
	mappers = []map[string]interface{}{}

	value := reflect.ValueOf(obj)
	for i := 0; i < value.Len(); i = i + 1 {
		item := value.Index(i)
		mappers = append(mappers, ObjToMap(item.Interface(), useTag))
	}
	return
}

func MapStringToObj(obj interface{}, mapper map[string]string, useTag string) (fields []string, errInfo string) {
	mapper2 := map[string]interface{}{}
	for k, v := range mapper {
		mapper2[k] = v
	}
	return MapToObj(obj, mapper2, useTag)
}

func setValue(v reflect.Value, value interface{}) (err error) {
	switch v.Kind() {
	case reflect.Bool:
		b, ok := Interface2Bool(value, false)
		if ok {
			v.SetBool(b)
		}
	case reflect.String:
		s, ok := Interface2String(value, false)
		if ok {
			v.SetString(s)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, ok := Interface2Int(value, false)
		if ok {
			v.SetInt(i)
		}
	case reflect.Float32, reflect.Float64:
		f, ok := Interface2Float(value, false)
		if ok {
			v.SetFloat(f)
		}
	//case reflect.Slice:
	//	vv := reflect.ValueOf(value)
	//	if vv.Kind() == reflect.Array || v.Kind() == reflect.Slice {
	//		l := vv.Len()
	//		newV := reflect.MakeSlice(v.Type(), l, l)
	//		for i := 0; i < l; i++ {
	//			setValue(newV.Index(i), vv.Index(i).Interface())
	//		}
	//		v.Set(newV)
	//	}
	default:
		// 非基本类型
		defer func() {
			e := recover()
			if e != nil {
				err = fmt.Errorf("%s %v", v.Type().String(), e)
			}
		}()
		v.Set(reflect.ValueOf(value))
		break
	}
	return
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

func Interface2StringWithType(value interface{}, strict bool) (v string, ok bool) {
	switch value.(type) {
	case string:
		v, ok = "string:" + value.(string), true
	case []uint8:
		v, ok = "[]uint8:" + string(value.([]uint8)), true
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
		v, ok = "int:" + strconv.FormatInt(i, 10), true
	case float64, float32:
		f, _ := Interface2Float(value, true)
		v, ok = "float:" + strconv.FormatFloat(f, 'f', -1, 64), true
	case bool:
		v, ok = "bool:" + strconv.FormatBool(value.(bool)), true
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

		if !field.CanInterface() {
			continue
		}

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
