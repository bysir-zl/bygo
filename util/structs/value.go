package structs

import (
	"strconv"
	"reflect"
	"strings"
	"errors"
)

func setValue(v reflect.Value, value interface{}) (err error) {
	if !v.CanSet() {
		err = errors.New("can't set")
		return
	}
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

	case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint:
		i, ok := Interface2UInt(value, false)
		if ok {
			v.SetUint(i)
		}
	case reflect.Float32, reflect.Float64:
		f, ok := Interface2Float(value, false)
		if ok {
			v.SetFloat(f)
		}

	default:
		// 非基本类型
		vx := reflect.ValueOf(value)
		if vx.Kind() == v.Kind() {
			v.Set(vx)
		}
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
	case int16:
		v, ok = int64(value.(int16)), true
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

func Interface2UInt(value interface{}, strict bool) (v uint64, ok bool) {
	switch value.(type) {
	case uint:
		v, ok = uint64(value.(uint)), true
	case uint8:
		v, ok = uint64(value.(uint8)), true
	case uint16:
		v, ok = uint64(value.(uint16)), true
	case uint32:
		v, ok = uint64(value.(uint32)), true
	case uint64:
		v, ok = uint64(value.(uint64)), true
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
		i, err := strconv.ParseUint(s, 10, 64)
		if err == nil {
			v, ok = i, true
		}
	case int, int8, int16, int32, int64:
		s, _ := Interface2Int(value, true)
		v, ok = uint64(s), true
	case float32, float64:
		f, _ := Interface2Float(value, true)
		v, ok = uint64(f), true
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
	case int8, int, int32, int64,
	uint8, uint, uint32, uint64:
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
	case uint, uint8, uint32, uint64:
		i, ok := Interface2UInt(value, true)
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
	case uint64, uint8, uint32, uint:
		i, _ := Interface2UInt(value, true)
		v, ok = strconv.FormatUint(i, 10), true
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
		v, ok = "string:"+value.(string), true
	case []uint8:
		v, ok = "[]uint8:"+string(value.([]uint8)), true
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
		v, ok = "int:"+strconv.FormatInt(i, 10), true
	case uint64, uint8, uint32, uint:
		i, _ := Interface2UInt(value, true)
		v, ok = "int:"+strconv.FormatUint(i, 10), true
	case float64, float32:
		f, _ := Interface2Float(value, true)
		v, ok = "float:"+strconv.FormatFloat(f, 'f', -1, 64), true
	case bool:
		v, ok = "bool:"+strconv.FormatBool(value.(bool)), true
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
