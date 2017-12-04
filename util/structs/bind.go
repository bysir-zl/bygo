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

func MapString2Struct(mapper map[string]string, obj interface{}, useTag string) (err error) {
	mapper2 := map[string]interface{}{}
	for k, v := range mapper {
		mapper2[k] = v
	}
	return Map2Struct(mapper2, obj, useTag)
}

// 将map将转换为struct，支持组合与嵌套
func Map2Struct(m map[string]interface{}, v interface{}, useTag string) (error) {
	return Map2StructValue(m, reflect.ValueOf(v), useTag)
}

func Map2StructValue(m map[string]interface{}, v reflect.Value, useTag string) (error) {
	pointer := indirect(v, false)
	typer := pointer.Type()

	fieldNum := pointer.NumField()

	for i := 0; i < fieldNum; i++ {
		field := pointer.Field(i)
		fieldT := typer.Field(i)
		// 如果是匿名 则需要扁平化
		if fieldT.Anonymous {
			e := Map2StructValue(m, field, useTag)
			if e != nil {
				return e
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

		if _, ok := m[fieldName]; !ok {
			continue
		}

		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.Struct:
			// 嵌套结构体
			if m2, ok := m[fieldName].(map[string]interface{}); ok {
				err := Map2StructValue(m2, field, useTag)
				if err != nil {
					return err
				}
			}
		case reflect.Slice:
			// 数组
			vv := reflect.ValueOf(m[fieldName])
			if vv.Kind() == reflect.Slice {
				l := vv.Len()

				newV := reflect.MakeSlice(field.Type(), l, l)
				for i := 0; i < l; i++ {
					// 这里只支持设置普通类型, 不支持结构体
					err := setValue(newV.Index(i), vv.Index(i).Interface())
					if err != nil {
						return err
					}
				}
				err := setValue(field, newV.Interface())
				if err != nil {
					return err
				}
			}
		default:
			setValue(field, m[fieldName])
		}
	}

	return nil
}
