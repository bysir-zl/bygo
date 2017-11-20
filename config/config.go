package config

import (
	"github.com/go-ini/ini"
	"os"
)

var f *ini.File

func GetInt(key string, section string) (value int) {
	value, _ = f.Section(section).Key(key).Int()
	return
}
func GetString(key string, section string) (value string) {
	value = f.Section(section).Key(key).String()
	return
}

func GetBool(key string, section string) (value bool) {
	value, _ = f.Section(section).Key(key).Bool()
	return
}

func Keys(section string) []string {
	ks := f.Section(section).Keys()
	r := make([]string, len(ks))
	for i, l := range ks {
		r[i] = l.Name()
	}
	return r
}

// filePath "config/app.ini"
func Load(filePath string) {
	var err error
	// 向上层查找文件
	// 在项目的任何地方运行(test时)都能加载到配置文件
	var file *os.File
	for i := 0; i < 10; i++ {
		f, err := os.Open(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				filePath = "../" + filePath
			} else {
				panic(err)
			}
		} else {
			file = f
			break
		}
	}
	if file == nil {
		panic("can't find config file")
	}
	f, err = ini.Load(file)
	if err != nil {
		panic(err)
	}
	f.BlockMode = false
}
