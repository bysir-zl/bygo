package config

import "bygo/db"

var DbConfigs map[string]db.DbConfig = map[string]db.DbConfig{}

func init() {
    if !Debug {
        DbConfigs["default"] = db.DbConfig{
            Driver:"mysql",
            Host:"localhost",
            Port:3306,
            Name:"password",
            User:"root",
            Password:"zhangliang",
        }
    } else {
        DbConfigs["default"] = db.DbConfig{
            Driver:"mysql",
            Host:"localhost",
            Port:3306,
            Name:"password",
            User:"root",
            Password:"",
        }
    }
}
