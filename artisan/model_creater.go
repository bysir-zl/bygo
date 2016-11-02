package artisan

import (
	"github.com/bysir-zl/bygo/config"
	"github.com/bysir-zl/bygo/db"
	"lib.com/deepzz0/go-com/log"
	"github.com/bysir-zl/bygo/util"
	"fmt"
	"os"
	"strings"
)

func CreateModelFile(tableName string) {
	conf, _ := config.LoadConfigFromFile("G:/go_project/src/project.com/bysir-zl/kuaifa-api/config/config.json")
	dbf := db.NewDbFactory(conf.DbConfigs)
	lis, err := dbf.Model(nil).Query("show columns from " + tableName)
	if err != nil {
		log.Warn(err)
	}

	modelName := string(util.SheXing2TuoFeng([]byte(tableName)))

	fileContent := "package model \r\n\r\n" + "type " + modelName + "Model struct {\r\n"
	fileContent = fileContent + "	Table string `db:\"" + tableName + "\" json:\"-\"`\r\n"
	fileContent = fileContent + "	Connect string `db:\"default\" json:\"-\"`\r\n\r\n"
	for _, v := range lis {
		pk := ""
		if string(v["Key"].([]uint8)) == "PRI" {
			pk = ` pk:"" `
			if string(v["Extra"].([]uint8)) == "auto_increment" {
				pk = ` pk:"auto" `
			}
		}

		t := string(v["Type"].([]uint8))
		key := string(v["Field"].([]uint8))
		tyo := "string"
		if strings.Contains(t, "int") {
			tyo = "int"
		} else if strings.Contains(t, "timestamp") {
			tyo = "int"
		}

		fieldName := string(util.SheXing2TuoFeng([]byte(key)))
		line := "	" + fieldName + " " + tyo + fmt.Sprintf(" `" + `name:"%s" ` + pk + `  json:"%s"` + "`", key, key) + "\r\n"
		fileContent = fileContent + line
	}
	fileContent = fileContent + "}"

	file, err := os.Create(tableName + "_model.go")
	if err != nil {
		log.Warn(err)
	}
	file.WriteString(fileContent)
}


