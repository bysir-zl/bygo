package artisan

import (
	"os"
	"io/ioutil"
	"github.com/deepzz0/go-com/log"
	"strings"
)

type params struct {
	name             string
	in               string
	description      string
	required         bool
	types            string `json:"type"`
	items            struct {
				 types    string `json:"type"`
				 enum     []string
				 defaults string `json:"default"`
			 }
	defaults         string `json:"default"`
	collectionFormat string
}
type bpi struct {
	tags        []string
	summary     string
	description string
	operationId string
	parameters  []params
}
type api struct {
	methods map[string]bpi
}

func parseFile() {
	fileName := "./swagger_test/index_controller.go"
	file, err := os.OpenFile(fileName, os.O_APPEND, 0666)
	if err != nil {
		log.Warn(err)
	}
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		log.Warn(err)
	}
	// 遍历每一行
	// 是否包含 // @
	str := string(bs)
	str=strings.Replace(str,"\r","",-1)
	ss := strings.Split(str, "\n")

	inOneApi := false

	var apiStrings []string = []string{}
	var apiString string
	for _, s := range ss {
		if strings.Contains(s, "// @") {
			if inOneApi {
				apiString =apiString+ strings.Replace(s,"// @","",-1)+ "\n"
			} else {

				inOneApi = true
				apiString =apiString+ strings.Replace(s,"// @","",-1)+ "\n"
			}
		} else if len(strings.Trim(strings.Trim(s, " "), "/")) > 1 {
			if inOneApi{
				apiStrings = append(apiStrings,apiString)
				apiString = ""
				inOneApi = false
			}
		}
	}

	var a api
	for _,item :=range apiStrings{
		item = strings.Replace(item," ","",-1)
		item = strings.Replace(item,";\n",";",-1)
		item = strings.Replace(item,":\n",":",-1)
		row:=strings.Split(item,"\n")

		a  = api{
			methods:map[string]bpi{},
		}
		desc:=row[0]
		url:=getRowString(row,"router")

		a.methods[url]=bpi{
			description:desc,
		}
	}

	log.Print(a.methods)
}

func getRowString(row []string,number string) string {
	for _,st:=range row{
		if strings.Index(st,number+":")==0{
			return strings.Split(st,":")[1]
		}
	}
	return ""
}