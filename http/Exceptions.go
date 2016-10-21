package http

type Exceptions struct {
    Err error //详细
    Code int //错误码
    Message string //简述
}