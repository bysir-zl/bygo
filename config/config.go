package config

import (
    "github.com/bysir-zl/bygo/db"
    "os"
    "io/ioutil"
    "github.com/widuu/gojson"
    "strconv"
    "github.com/bysir-zl/bygo/util"
)


type Config struct {
    Evn        string
    Debug      bool
    CacheDrive string
    DbConfigs  map[string]db.DbConfig
    RedisHost  string
}

func LoadConfigFromFile(filePath string) (config Config ,err error) {
    file, err := os.Open(filePath)
    if err != nil {
        return
    }
    bs, err := ioutil.ReadAll(file)
    if err != nil {
        return
    }
    jsonString := string(bs)

    config = Config{}
    config.Evn = gojson.Json(jsonString).Get("evn").Tostring()
    config.Debug, _ = strconv.ParseBool(gojson.Json(jsonString).Get("app").Get("debug").Tostring())
    config.CacheDrive = gojson.Json(jsonString).Get("cache").Get("cache_driver").Tostring()
    config.RedisHost = gojson.Json(jsonString).Get("redis").Get("host").Tostring()

    ds := gojson.Json(jsonString).Get("database").Getdata()
    dbConfigs := map[string]db.DbConfig{}
    for name, conf := range ds {
        c := db.DbConfig{}
        util.MapToObj(&c, conf.(map[string]interface{}), "json")
        dbConfigs[name] = c
    }
    config.DbConfigs = dbConfigs

    return
}