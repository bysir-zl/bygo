package artisan

import (
	"github.com/atotto/clipboard"
	"github.com/bysir-zl/bygo/util"
	"regexp"
	"errors"
)

//
func Xml2Go() (err error) {
	res, err := clipboard.ReadAll()
	if err != nil {
		return
	}

	tr := "type GoStruct struct {\n"

	root := ""
	{
		r, _ := regexp.Compile(`<(.*?)>`)
		result := r.FindStringSubmatch(res)
		if len(result) != 2 {
			return errors.New("can't get root")
		}
		root = result[1]
	}

	tr += "XMLName        struct{} `xml:\"" + root + "\"`\n"
	{
		r, _ := regexp.Compile(`<(.*?)><!\[CDATA\[(.*?)\]\]></.*>`)
		result := r.FindAllStringSubmatch(res, -1)
		for _, v := range result {
			field := v[1]
			fieldName := string(util.SheXing2TuoFeng([]byte(field)))
			tr += "    " + fieldName + " " + "string" + " `xml:\"" + field + "\"`" + "\n"
		}
	}

	tr += "}"

	err = clipboard.WriteAll(tr)
	return
}
