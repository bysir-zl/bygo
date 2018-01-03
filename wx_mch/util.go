package wx_mch

import (
	"reflect"
	"fmt"
	"strings"
	"github.com/bysir-zl/bygo/util"
	"github.com/bysir-zl/bygo/util/encoder"
)

// 签名
func SignData(src interface{}, key string) (sign string) {
	kv := wxParseSignFields(src)
	s := kv.EncodeString()
	if len(s) > 0 {
		s += "&"
	}
	s += "key=" + key
	return encoder.Md5String(s)
}

func wxParseSignFields(src interface{}) util.OrderKV {
	values := util.OrderKV{}
	t := reflect.TypeOf(src)
	v := reflect.ValueOf(src)
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		if tf.Tag.Get("sign") != "true" {
			continue
		}
		tv := v.Field(i)
		if !tv.IsValid() {
			continue
		}
		sv := fmt.Sprintf("%v", tv.Interface())
		if sv == "" {
			continue
		}
		name := ""
		if xn := tf.Tag.Get("xml"); xn != "" {
			name = strings.Split(xn, ",")[0]
		} else if xn = tf.Tag.Get("json"); xn != "" {
			name = strings.Split(xn, ",")[0]
		} else {
			continue
		}
		values.Add(name, sv)
	}
	return values
}
