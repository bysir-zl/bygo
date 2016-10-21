package db

import (
    "fmt"
)

type DbConfig struct {
    Driver   string
    //
    Host     string
    //端口
    Port     int
    //用户名
    User     string
    //密码
    Password string
    //数据库名name
    Name     string
}

func (p DbConfig) String() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s~%s", p.User, p.Password, p.Host, p.Port, p.Name, p.Driver);
}
func (p DbConfig) GetSqlOpenString() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", p.User, p.Password, p.Host, p.Port, p.Name);
}
