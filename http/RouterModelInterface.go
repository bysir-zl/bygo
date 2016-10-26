package http

import (
    "github.com/bysir-zl/bygo/db"
    "github.com/bysir-zl/bygo/bean"
)

type  RouterModelInterface interface {
    //c 出于规范的考虑,不能实现自动注入,请自己从Container中取所需值
    //rs 如果你返回的ResponseData.Exception.Code!=0,则代表你直接拒绝了此次访问,bygo将直接将ResponseData输出,
    //mr 你还必须指定规则,比如 [在更新的时候,只能更新指定字段],[在查询的时候,指定返回什么字段,指定只能以什么为条件,以什么Order],[在添加的时候,指定添加那些字段]
    SelectBefore(c Context) (rs ResponseData, mr db.ModelRules)
    //rs默认返回值,[Select返回*slice(model),Insert返回*model,Delete返回个数,Update返回*model];
    SelectAfter(c Context, data interface{}, page bean.Page) (ResponseData)

    UpdateBefore(c Context) (rs ResponseData, mr db.ModelRules)
    UpdateAfter(c Context, data interface{}) (ResponseData)

    InsertBefore(c Context) (rs ResponseData, mr db.ModelRules)
    // r:要返回的自定义内容,如没有则返回默认响应; data:插入成功后的数据,*model类型
    InsertAfter(c Context, data interface{}) (ResponseData)

    DeleteBefore(c Context) (rs ResponseData, mr db.ModelRules)
    DeleteAfter(c Context, count int64) (ResponseData)
}