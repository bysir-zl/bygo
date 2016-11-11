package artisan

import (
	"os"
	"io/ioutil"
	"github.com/deepzz0/go-com/log"
	"strings"
	"strconv"
	"github.com/bysir-zl/bygo/util"
	"encoding/json"
)

type params struct {
	Name             string `json:"name"`
	In               string `json:"in"`
	Description      string `json:"description"`
	Required         bool `json:"required"`
	Types            string `json:"type"`
	Items            struct {
				 Types    string `json:"type,omitempty"`
				 Enum     []string `json:"enum,omitempty"`
				 Defaults string `json:"default,omitempty"`
			 }   `json:"items"`
	Defaults         interface{} `json:"default,omitempty"`
	CollectionFormat string `json:"collectionFormat,omitempty"`
}

type Schema struct {
	Ref string `json:"&ref"`
	// todo 暂时不支持这种格式
	//Types string `json:"type,omitempty"`
	//Items struct {
	//	      Ref string `json:"$ref,omitempty"`
	//      } `json:"items,omitempty"`
}

type response struct {
	Description string `json:"description"`
	Schema      Schema  `json:"schema,omitempty"`
}
type bpi struct {
	Tags        []string `json:"tags,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description"`
	OperationId string `json:"operation_id"`
	Parameters  []params `json:"parameters"`
	Responses   map[string]response `json:"responses"`
}

type router map[string]map[string]bpi

type info struct {
	Description string `json:"description"`
	Version     string `json:"version"`
	Title       string `json:"title"`
	Contact     struct {
			    Email string `json:"email,omitempty"`
		    } `json:"contact"`

	// todo licence
	License     struct {
			    Name string `json:"name,omitempty"`
			    Url  string `json:"url,omitempty"`
		    }  `json:"license"`
}

type Propertie struct {
	Types       string `json:"type"`
	Ref         string `json:"$ref,omitempty"`
	Default     string `json:"default"`
	Description string `json:"description"`
}
type Definition struct {
	Types      string `json:"type"`
	Properties map[string]Propertie `json:"properties"`
}

type swagger struct {
	Info        info `json:"info"`
	Swagger     string `json:"swagger"`
	Host        string `json:"host,omitempty"`
	BasePath    string `json:"basePath,omitempty"`
	Schemes     []string  `json:"schemes,omitempty"`
	Paths       router `json:"paths"`
	Definitions map[string]Definition `json:"definitions"`
}

type swaggerString map[string][]string

func getAllSwaggerString(root string) (sw swaggerString) {
	sw = swaggerString{}
	paths, err := util.WalkDir(root, []string{"go", "proto"})
	if err != nil {
		log.Warn(err)
		return
	}

	for _, fileName := range paths {
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
		str = strings.Replace(str, "\r", "", -1)
		ss := strings.Split(str, "\n")
		// 添加最后一个空行,用于结束匹配
		ss = append(ss, "")

		types := ""

		var apiString string
		for _, s := range ss {
			if strings.Index(s, "// @") == 0 {
				if types != "" {
					// 非第一行 ,去除空格
					s = strings.Replace(s, "// @", "", -1)
					//s = strings.Replace(s, " ", "", -1)
					apiString = apiString + s + "\n"
				} else {
					// 是第一行
					if strings.Contains(s, "@API ") {
						types = "API"
					} else if strings.Contains(s, "@BASE ") {
						types = "BASE"
					} else if strings.Contains(s, "@INFO ") {
						types = "INFO"
					} else if strings.Contains(s, "@DEF ") {
						types = "DEF"
					} else {
						continue
					}
					// 将第一行的type标记去掉
					s = strings.Replace(s, "// @" + types + " ", "", -1)
					apiString = apiString + s + "\n"
				}
			} else if strings.Index(s, "//") != 0 {
				if types != "" {
					apiString = strings.Replace(apiString, ";\n", ";", -1)
					apiString = strings.Replace(apiString, ":\n", ":", -1)
					// 去掉最后一个换行
					apiString = apiString[:len(apiString) - 1]
					if sw[types] == nil {
						sw[types] = []string{apiString}
					} else {
						sw[types] = append(sw[types], apiString)
					}

					apiString = ""
					types = ""
				}
			}
		}
	}

	return
}

func S(root string, output string) (err error) {
	sw := getAllSwaggerString(root)
	base := map[string]string{}
	rou := router{}
	inf := map[string]string{}
	defini := map[string]Definition{}
	if sw["BASE"] != nil {
		base = parseBase(sw["BASE"])
	}
	if sw["API"] != nil {
		rou = parsePath(sw["API"], base)
	}
	if sw["INFO"] != nil {
		inf = parseInfo(sw["INFO"])
	}
	if sw["DEF"] != nil {
		defini = parseDef(sw["DEF"])
	}

	title := inf["title"]
	desc := inf["desc"]
	version := inf["version"]
	email := inf["email"]
	host := inf["host"]
	basePath := inf["basePath"]

	swagger := swagger{
		Swagger:"2.0",
		Info:info{
			Description:desc,
			Title:title,
			Version:version,
		},
		Paths:rou,
		Host:host,
		BasePath:basePath,
		Definitions:defini,
	}
	swagger.Info.Contact.Email = email

	bs, _ := json.MarshalIndent(&swagger, "", "    ")

	file, e := os.Create(output)
	if e != nil {
		file, e = os.Open(output)
		if e != nil {
			err = e
			return
		}
	}
	defer func(file *os.File) {
		file.Close()
	}(file)

	s := string(bs)
	s = strings.Replace(s, `\u0026`, "$", -1)

	file.WriteString(s)
	return
}

func parseDef(ss []string) map[string]Definition {
	defs := map[string]Definition{}
	for _, item := range ss {
		row := strings.Split(item, "\n")
		nameAndType := strings.Split(row[0], ",")
		name := nameAndType[0]
		types := "object"
		if len(nameAndType) > 1 {
			types = nameAndType[1]
		}
		if types == "object" {
			ps := map[string]Propertie{}

			psList := strings.Split(row[1], ";")
			for _, s := range psList {
				if len(strings.Split(s, ":")) < 2 {
					continue
				}
				pname := strings.Split(s, ":")[0]
				s = strings.Split(s, ":")[1]
				desc := ""
				t := "string"
				defau := ""
				args := strings.Split(s, ",")
				if len(args) > 0 {
					desc = args[0]
				}
				if len(args) > 1 {
					t = args[1]
				}
				if len(args) > 2 {
					defau = args[2]
				}
				if t == "ref" {
					ps[pname] = Propertie{
						Ref:"#/definitions/" + defau,
					}
				} else {
					if defau == "" {
						defau = desc
					}
					ps[pname] = Propertie{
						Default:defau,
						Description:desc,
						Types:t,
					}
				}
			}

			defs[name] = Definition{
				Properties:ps,
				Types:"object",
			}
		}

	}

	return defs
}
func parseBase(ss []string) map[string]string {
	base := map[string]string{}
	for _, item := range ss {
		row := strings.Split(item, "\n")

		name := row[0]
		text := strings.Join(row[1:], "\n")
		base[name] = text
	}
	return base
}

func parsePath(apis []string, base map[string]string) router {
	rou := router{}

	for _, api := range apis {
		// 替换BASE

		if strings.Contains(api, "BASE:") {
			replaced := strings.Split(api, "BASE:")[1]
			replaced = strings.Split(strings.Split(replaced, "\n")[0], ";")[0]
			api = strings.Replace(api, "BASE:" + replaced, base[replaced], -1)
		}

		row := strings.Split(api, "\n")

		// desc
		desc := row[0]
		// router
		router := getRowString(row, "router")
		if len(router) == 0 {
			log.Warn(row)
			return rou
		}
		routers := strings.Split(router, ",")
		url := routers[0]
		method := routers[1]
		tags := strings.Split(routers[2], "|")
		operationId := routers[3]

		bpii := bpi{
			Description:desc,
			Tags:tags,
			OperationId:operationId,
		}

		// parameters
		para := getRowString(row, "parameters")

		if len(para) != 0 {
			paras := strings.Split(para, ";")
			paList := []params{}
			for _, p := range paras {
				if !strings.Contains(p, ":") || !strings.Contains(p, ",") {
					continue
				}
				pos:=strings.Index(p,":")
				name := p[:pos]
				p := p[pos+1:]
				ps := strings.Split(p, ",")
				if len(ps) < 3 {
					continue
				}
				desc := ps[0]
				types := ps[1]
				in := ps[2]
				var defaults interface{}
				if len(ps) > 3 {
					switch types {
					case "boolean":
						defaults, _ = strconv.ParseBool(ps[3])
					case "int":
						defaults, _ = strconv.ParseInt(ps[3], 10, 64)
					case "string":
						defaults = ps[3]
					default:
						defaults = ps[3]
					}
				}
				required := false
				if len(ps) > 4 {
					required, _ = strconv.ParseBool(ps[4])
				}
				paList = append(paList, params{
					//CollectionFormat:"multi",
					Defaults:defaults,
					Description:desc,
					In:in,
					Name:name,
					Types:types,
					Required:required,
				})
			}
			bpii.Parameters = paList
		}

		// responses
		res := map[string]response{}
		respon := getRowString(row, "responses")

		if len(respon) != 0 {
			respons := strings.Split(respon, ";")
			for _, r := range respons {
				pos := strings.Index(r, ":")

				name := r[:pos]
				r = r[pos + 1:]
				rs := strings.Split(r, ",")
				desc := rs[0]
				ref := ""
				if len(rs) > 1 {
					ref = "#/definitions/" + rs[1]
				}
				res[name] = response{
					Description:desc,
					Schema:Schema{
						Ref:ref,
					},
				}
			}
			bpii.Responses = res
		}
		// end

		rou[url] = map[string]bpi{}
		rou[url][method] = bpii
	}

	return rou
}

func parseInfo(ss []string) map[string]string {
	s := ss[0]
	info := map[string]string{}
	row := strings.Split(s, "\n")
	title := row[0]
	info["title"] = title
	info["desc"] = getRowString(row, "desc")
	info["version"] = getRowString(row, "version")
	info["email"] = getRowString(row, "email")
	info["host"] = getRowString(row, "host")
	info["basePath"] = getRowString(row, "basePath")
	info["license"] = getRowString(row, "license")

	return info
}

func Swagger(path, output string) {
	err := S(path, output)
	if err != nil {
		panic(err)
		return
	}
	log.Info("SUCCESS")
}

func getRowString(row []string, number string) string {
	for _, st := range row {
		if strings.Index(st, number) == 0 {
			pos := strings.Index(st, ":")
			s := st[pos + 1:]
			if s == "-" {
				return ""
			}
			return s
		}
	}
	return ""
}