package artisan
//
//import (
//	"fmt"
//	"github.com/bysir-zl/bygo/util"
//	"os"
//	"strings"
//	"github.com/bysir-zl/orm"
//	"github.com/bysir-zl/bygo/log"
//)
//
//func CreateModelFile(tableName string) {
//	config := map[string]orm.Connect{
//		"default":{
//			Driver:   "mysql",
//			Host:     "localhost",
//			Port:     3306,
//			Name:     "anyminisdk",
//			User:     "root",
//			Password: "root",
//		},
//	}
//	dbf := orm.New(config)
//	has, lis, err := dbf.QuerySql("show columns from " + tableName)
//	if err != nil {
//		log.Error("CreateModel",err)
//		return
//	}
//	if !has {
//		log.Info("CreateModel","can not read table")
//		return
//	}
//
//	modelName := string(util.SheXing2TuoFeng([]byte(tableName)))
//
//	fileContent := "package model \r\n\r\n" + "type " + modelName + "Model struct {\r\n"
//	fileContent = fileContent + "	orm string `table:\"" + tableName + "\" connect:\"default\" json:\"-\"`\r\n\r\n"
//	for _, v := range lis {
//		pk := ""
//		if string(v["Key"].([]uint8)) == "PRI" {
//			pk = ` pk:"" `
//			if string(v["Extra"].([]uint8)) == "auto_increment" {
//				pk = ` pk:"auto" `
//			}
//		}
//
//		t := string(v["Type"].([]uint8))
//		key := string(v["Field"].([]uint8))
//		tyo := "string"
//		if strings.Contains(t, "int") {
//			tyo = "int"
//		} else if strings.Contains(t, "timestamp") {
//			tyo = "int"
//		} else if strings.Contains(t, "bool") {
//			tyo = "bool"
//		}
//
//		fieldName := string(util.SheXing2TuoFeng([]byte(key)))
//		line := "	" + fieldName + " " + tyo + fmt.Sprintf(" `" + `db:"%s" ` + pk + `  json:"%s"` + "`", key, key) + "\r\n"
//		fileContent = fileContent + line
//	}
//	fileContent = fileContent + "}"
//
//	file, err := os.Create(tableName + "_model.go")
//	if err != nil {
//		log.Warn(err)
//	}
//	file.WriteString(fileContent)
//	log.Print("create success")
//
//}
