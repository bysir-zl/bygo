package http

import (
    "reflect"
    "strings"
    "errors"
    "strconv"
    "github.com/bysir-zl/bygo/util"
    "github.com/bysir-zl/bygo/bean"
)

//用于处理Router中的Model
type RouterModelHandler struct {
    model  RouterModelInterface
    method string
}

func (p *RouterModelHandler) Handle(container Context) ResponseData {
    method := p.method

    if method == "Update" {
        return p.handleUpdate(container)
    }
    if method == "Select" {
        return p.handleSelect(container)
    }
    if method == "Delete" {
        return p.handleDelete(container)
    }
    if method == "Insert" {
        return p.handleInsert(container)
    }

    if (util.ItemInArray(method, []string{
        "BeforeUpdate",
        "AfterUpdate",
        "BeforeSelete",
        "AfterSelete",
        "BeforeDelete",
        "AfterDelete",
        "BeforeInsert",
        "AfterInsert",
    })) {
        return NewRespDataError(403, errors.New("forbidden"))
    }

    model := p.model
    fv := reflect.ValueOf(model)
    me := fv.MethodByName(method)

    if !me.IsValid() {
        return NewRespDataError(404, errors.New("the method '" + method + "' is undefined " +
            "in model '" + reflect.TypeOf(model).String() + "'!"))
    }

    params, err := container.GetFuncParams(me);
    if err != nil {
        return NewRespDataError(500, err)
    }

    response := (me.Call(params)[0].Interface()).(ResponseData)
    return response;

}

// 处理update方法
func (p *RouterModelHandler) handleUpdate(container Context) ResponseData {
    model := p.model

    response, modelRules := model.UpdateBefore(container)

    //如果返回不为空,则视为拒绝访问,直接返回错误信息
    if response.Code != 0 {
        return response;
    }
    dbFactory := container.DbFactory
    request := container.Request
    fields := request.Input.SetToObjWithDenyFilter(model, "Where")

    err := modelRules.CheckFields(fields)
    if err != nil {
        return NewRespDataError(403, err)
    }

    factory := dbFactory.Model(model);

    // where
    whereString := request.Input.Get("Where") // "Id:5"
    where, err := modelRules.CheckWhereString(whereString, true);
    if err != nil {
        return NewRespDataError(403, err)
    }
    for wk, wv := range where {
        factory = factory.Where(wk, wv...)
    }

    count, err :=
        factory.Update()
    if err != nil {
        return NewRespDataError(500, err)
    }

    //默认返回值
    defaultResponse := NewRespDataJson(200, bean.ApiData{Code:200})
    if count == 0 {
        defaultResponse = NewRespDataJson(200, bean.ApiData{Code:200, Msg:"执行成功,但没有更新"})
    }

    responseData2 := model.UpdateAfter(container, model);
    if responseData2.Code != 0 {

        return responseData2
    }
    return defaultResponse
}

func (p *RouterModelHandler) handleSelect(container Context) ResponseData {
    model := p.model
    response, modelRules := model.SelectBefore(container)

    //如果返回不为空,则视为拒绝访问,直接返回错误信息
    if response.Code != 0 {
        return response;
    }
    dbFactory := container.DbFactory
    request := container.Request
    whereString := request.Input.Get("Where") // "Id:5"

    //生成一个Slice用于存储返回列表
    models := reflect.New(reflect.SliceOf(reflect.TypeOf(model).Elem())).Interface()
    factory := dbFactory.Model(models)

    where, err := modelRules.CheckWhereString(whereString, false)
    if err != nil {
        return NewRespDataError(403, err)
    }
    for condition, values := range where {
        factory = factory.Where(condition, values...)
    }

    // 解析Order
    orderString := request.Input.Get("Order") // "Id:DESC"
    order, err := modelRules.CheckOrderString(orderString)
    if err != nil {
        return NewRespDataError(403, err)
    }
    for ok, ov := range order {
        factory = factory.OrderBy(ok, ov)
    }

    //设置查询字段
    field := modelRules.GetFieldsList()
    if (field[0] != "*") {
        factory.Field(field...)
    } else {

    }

    // 处理分页
    page, _ := strconv.ParseInt(request.Input.Get("Page"), 10, 64);
    pageSize, _ := strconv.ParseInt(request.Input.Get("PageSize"), 10, 64);
    if pageSize == 0 {
        pageSize = 20
    }

    var pageData bean.Page
    var err2 error
    var apiReturn interface{};

    // 不分页
    if modelRules.GetPageRules().WithOutPage {
        err2 = factory.Get()
        apiReturn = bean.ApiData{Data: models, Code:200}
    } else {
        // 不计总数
        if modelRules.GetPageRules().WithOutTotal {
            pageData, err2 =
                factory.PageWithOutTotal(int(page), int(pageSize))
        } else {
            pageData, err2 =
                factory.Page(int(page), int(pageSize))
        }
        apiReturn = bean.ApiDataWithPage{Page:pageData, Data: models, Code:200}
    }
    if err2 != nil {
        return NewRespDataError(500, err2)
    }

    defaultResponse := NewRespDataJson(200, apiReturn)

    responseData2 := model.SelectAfter(container, models, pageData);
    if responseData2.Code != 0 {
        return responseData2
    }
    return defaultResponse
}

func (p *RouterModelHandler) handleInsert(container Context) ResponseData {
    model := p.model
    response, modelRules := model.InsertBefore(container)
    util.EmptyObject(model);

    //如果返回不为空,则视为拒绝访问,直接返回错误信息
    if response.Code != 0 {
        return response;
    }
    dbFactory := container.DbFactory

    request := container.Request
    fields := request.Input.SetToObj(model, "json")

    factory := dbFactory.Model(model).Field(fields...);

    if in, msg := util.ArrayInArray(fields, modelRules.GetFieldsList()); !in {
        return NewRespDataError(403, errors.New("服务器拒绝插入" + msg + "字段"))
    }

    err := factory.Insert();
    if err != nil {
        return NewRespDataError(500, err)
    }

    defaultResponse := NewRespDataJson(200, bean.ApiData{Code:200, Data:model})

    responseData2 := model.InsertAfter(container, model)
    if responseData2.Code != 0 {
        return responseData2
    }
    return defaultResponse
}

func (p *RouterModelHandler) handleDelete(container Context) ResponseData {
    model := p.model

    response, modelRules := model.DeleteBefore(container)
    //如果返回不为空,则视为拒绝访问,直接返回错误信息
    if response.Code != 0 {
        return response;
    }
    dbFactory := container.DbFactory
    request := container.Request

    whereString := request.Input.Get("Where") // "Id:5"
    if whereString == "" {
        return NewRespDataError(403, errors.New("呵呵,你是不是忘了传Where=Id:1"))
    }

    factory := dbFactory.Model(model);
    whereRules := modelRules.GetWhereRules();

    //解析Where字符串
    //Where=Id:20,Key:zl,Sex=1
    where := map[string]string{};
    whereKvs := strings.Split(whereString, ",")
    for _, kvStr := range whereKvs {
        kv := strings.Split(kvStr, ":")
        if len(kv) != 2 || kv[0] == "" || kv[1] == "" {
            return NewRespDataError(403, errors.New("呵呵,你是不是传错了Where"))
        }

        // 没指定此字段
        if whereRules[kv[0]] == "" {
            return NewRespDataError(403, errors.New("服务器拒绝通过" + kv[0] + "字段删除"))
        }

        where[kv[0]] = kv[1];
    }

    for wk, wv := range where {
        wString := whereRules[wk];
        factory = factory.Where(wString, wv)
    }

    defaultResponse := NewRespDataJson(200, bean.ApiData{Code:200})
    count, err :=
        factory.
        Delete()
    if err != nil {
        return NewRespDataError(500, err)
    }
    if count == 0 {
        defaultResponse = NewRespDataJson(200, bean.ApiData{Code:200, Msg:"执行成功,但没有删除"})
    }

    responseData2 := model.DeleteAfter(container, count);
    if responseData2.Code != 0 {
        return responseData2
    }
    return defaultResponse
}

