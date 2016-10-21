package exception

import (
    "app/bean"
    "bygo/http"
)

// 将报错的Exception处理成Response返回。在这里你可以判断e.Code统一处理错误,比如上报code==500的错误
func Handler(c http.SessionContainer, e http.Exceptions) http.ResponseData {
    return http.NewRespDataJson(200, bean.FailReturn{Code:e.Code, Msg:e.Message})
}
