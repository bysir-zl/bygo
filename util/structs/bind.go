package structs

import (
	"reflect"
	"strings"
)

// 将结构体转换为map[string]string{}，仅支持组合和普通一层结构体，对第二层结构体会直接忽略
func Struct2MapString(v interface{}, useTag string) (map[string]string, error) {
	pointer := reflect.Indirect(reflect.ValueOf(v))
	typer := pointer.Type()

	fieldNum := pointer.NumField()
	m := map[string]string{}

	for i := 0; i < fieldNum; i++ {
		field := pointer.Field(i)

		if !field.CanInterface() {
			continue
		}
		fieldT := typer.Field(i)

		// 如果是匿名 则需要扁平化
		if fieldT.Anonymous {
			mn, err := Struct2MapString(field.Interface(), useTag)
			if err != nil {
				return nil, err
			}
			for k, v := range mn {
				m[k] = v
			}
			continue
		}

		fieldName := fieldT.Name
		if useTag != "" {
			fieldName = fieldT.Tag.Get(useTag)
			if fieldName == "" {
				// 如果指定了tag 但是tag为空，则不处理这个字段
				continue
			}

			fieldName = strings.Split(fieldName, ",")[0]
		}

		switch field.Kind() {
		case reflect.Interface, reflect.Ptr, reflect.Array, reflect.Map, reflect.Slice:
		default:
			m[fieldName], _ = Interface2String(field.Interface(), false)
		}

	}
	return m, nil
}

// 将map将转换为struct，支持组合与嵌套
func Map2Struct(m map[string]interface{}, v interface{}, useTag string) (error) {
	pointer := reflect.Indirect(reflect.ValueOf(v))
	typer := pointer.Type()

	fieldNum := pointer.NumField()

	for i := 0; i < fieldNum; i++ {
		field := pointer.Field(i)
		fieldT := typer.Field(i)
		// 如果是匿名 则需要扁平化
		if fieldT.Anonymous {
			// 本来开始是直接把field.interface甩进去的, 但是得到的field是!CanSet的,
			// 所以这里直接新建一个 直接赋值整个结构体
			value := reflect.New(field.Type())
			t := value.Interface()
			e := Map2Struct(m, t, useTag)
			if e != nil {
				return e
			}
			field.Set(value.Elem())
			continue
		}

		fieldName := fieldT.Name
		if useTag != "" {
			fieldName = fieldT.Tag.Get(useTag)
			if fieldName == "" {
				// 如果指定了tag 但是tag为空，则不处理这个字段
				continue
			}

			fieldName = strings.Split(fieldName, ",")[0]
		}

		if _, ok := m[fieldName]; !ok {
			continue
		}

		if !field.CanSet() {
			continue
		}
		setValue(field, m[fieldName])
	}

	return nil
}
