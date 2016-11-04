package http

import (
	"github.com/bysir-zl/bygo/bean"
	"fmt"
)

type Logs struct {
	Err     error  //详细
	Code    int    //错误码
	Message string //简述
}

func (p *Logs) ToResponseData() ResponseData {
	return NewRespDataJson(200, bean.ApiData{Code: p.Code, Msg: p.Message})
}
func (p Logs) String()string  {
	return fmt.Sprintf(" 'Code :%d , Msg : %s' ",p.Code,p.Message)
}
