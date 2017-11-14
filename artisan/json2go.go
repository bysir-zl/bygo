package artisan

import (
	"github.com/atotto/clipboard"
	"encoding/json"
	"reflect"
	"github.com/bysir-zl/bygo/util"
)

//
func Json2Go() (err error) {
	j, err := clipboard.ReadAll()
	if err != nil {
		return
	}

	tr := "type GoStruct struct {\n"

	js := map[string]interface{}{}

	err = json.Unmarshal([]byte(j), &js)
	if err != nil {
		return
	}

	for k, v := range js {
		fieldName := string(util.SheXing2TuoFeng([]byte(k)))
		tr += "    " + fieldName + " " + reflect.TypeOf(v).String() + " `json:\"" + k + "\"`" + "\n"
	}

	tr += "}"

	err = clipboard.WriteAll(tr)
	return
}
