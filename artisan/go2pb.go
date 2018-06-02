package artisan

import (
	"regexp"
	"strconv"
	"strings"
	"github.com/bysir-zl/bygo/util"
	"github.com/atotto/clipboard"
	"errors"
)

// go 结构体转pb结构体

/*type GetSiteListParams struct {
	Page      int    `json:"page" form:"page"`
	Size      int    `json:"size" form:"size"`
	Id        int64  `json:"id" form:"id"`
	Status    int8   `json:"status" form:"status"`
	Email     string `json:"email" form:"email"`
	Mobile    string `json:"mobile" form:"mobile"`
	ServiceId int64  `json:"service_id" form:"service_id"`
	TimeStart int64  `json:"time_start" form:"time_start"`
	TimeEnd   int64  `json:"time_end" form:"time_end"`
}*/

func Go2Pb() (err error) {
	j, err := clipboard.ReadAll()
	if err != nil {
		return
	}

	r, err := go2pb(j)
	if err != nil {
		return
	}
	err = clipboard.WriteAll(r)
	return
}

func go2pb(in string) (out string, err error) {
	result := strings.Split(in, "\n")

	typesMap := map[string]string{
		"int":  "int32",
		"int8": "int32",
	}

	message := "message "
	isEmpty := true

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
			isEmpty = false
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

	if isEmpty {
		err = errors.New("empty input")
	} else {
		err = nil
	}
	return message, err
}
