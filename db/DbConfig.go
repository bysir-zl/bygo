package db

import (
    "fmt"
)

type DbConfig struct {
    Driver   string `json:"driver"`
    //
    Host     string `json:"host"`
    //端口
    Port     int `json:"port"`
    //用户名
    User     string `json:"user"`
    //密码
    Password string `json:"password"`
    //数据库名name
    Name     string `json:"name"`
}

func (p DbConfig) String() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s~%s", p.User, p.Password, p.Host, p.Port, p.Name, p.Driver);
}
func (p DbConfig) GetSqlOpenString() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", p.User, p.Password, p.Host, p.Port, p.Name);
}
