package artisan

import (
	"testing"
	"strings"
	"github.com/bysir-zl/bygo/util"
	"strconv"
	"regexp"
)

func TestGo2Pb(t *testing.T) {
	res := `
type GetSiteListParams struct {
	Page      int
	Size      int
	Id        int64
	Status    int8
	Email     string
	Mobile    string
	ServiceId int64
	TimeStart int64
	TimeEnd   []*User
}
`
	result := strings.Split(res, "\n")

	typesMap := map[string]string{
		"int":  "int32",
		"int8": "int32",
	}

	message := "message "
	for i, vs := range result {
		if vs == "" {
			continue
		}

		// 去掉多余空格
		r, _ := regexp.Compile(`\s+`)
		vs = r.ReplaceAllString(vs, " ")

		vs = strings.Trim(vs, " ")
		v := strings.Split(vs, " ")
		if len(v) < 2 {
			continue
		}

		if v[0] == "type" {
			message += v[1] + " {\n"
			continue
		}
		fieldName := string(util.TuoFeng2SheXing([]byte(v[0])))
		types := v[1]

		if tm, ok := typesMap[types]; ok {
			types = tm
		}
		message += "    "

		// 判断数组
		if len(types) >= 2 && types[:2] == "[]" {
			types = types[2:]
			message += "repeated "
		}
		types = strings.Trim(types, "*")
		message += types + " " + fieldName + " = " + strconv.Itoa(i) + ";\n"
	}

	message += "}\n"

	t.Log(message)

}
