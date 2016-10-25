package http

import (
    "reflect"
    "strings"
    "errors"
    "github.com/bysir-zl/bygo/db"
    "github.com/bysir-zl/bygo/cache"
)


// 存储用于依赖注入的容器
type Container struct {
    DbFactory    db.DbFactory
    Cache        cache.CacheInterface

    OtherItemMap map[string]interface{}
}

func (p *Container)GetItemByClassName(name string) interface{} {
    return p.OtherItemMap[name]
}

func (p *Container)GetItemByClass(item interface{}) interface{} {
    va := reflect.ValueOf(item);
    return p.OtherItemMap[va.Type().String()]
}

// 向容器中注册一个Item
// eg : SetItem(App.User{}) ,这样在使用函数或者方法时使用func(u App.User)会自动注入当前App.User
func (p *Container)SetItem(item interface{}) {
    va := reflect.ValueOf(item);
    p.OtherItemMap[ va.Type().String()] = item
}

//从容器中获取参数类型
func (s *Container)GetFuncParams(fun reflect.Value) (data []reflect.Value, err error) {
    var params []reflect.Value = nil;

    rf := fun.Type().String();

    ps := strings.Split(rf, "(")[1];
    ps = strings.Split(ps, ")")[0];

    //如果有参数需要注入
    if len(ps) != 0 {
        ps = strings.Replace(ps, " ", "", -1)
        paras := strings.Split(ps, ",");

        params = make([]reflect.Value, len(paras))

        for index, cla := range paras {
            p := s.GetItemByClassName(cla);
            if (p == nil) {
                return nil, errors.New("container not has '" +
                    cla + "' item , please use container.SetItem() to set item")
            }
            params[index] = reflect.ValueOf(p)
        }
    }

    return params, nil;
}


// 存储用于依赖注入的容器
type SessionContainer struct {
    Response     *Response
    Request      *Request

    DbFactory    db.DbFactory
    Cache        cache.CacheInterface

    OtherItemMap map[string]interface{}
}

func (p *SessionContainer)GetItemByClass(item interface{}) interface{} {
    va := reflect.ValueOf(item)
    return p.GetItemByClassName(va.Type().String())
}

func (p *SessionContainer)GetItemByClassName(name string) interface{} {

    if (name == "db.DbFactory") {
        return p.DbFactory
    }
    if (name == "cache.CacheInterface") {
        return p.Cache
    }
    if (name == "*http.Request") {
        return p.Request
    }
    if (name == "*http.Response") {
        return p.Response
    }
    if (name == "http.SessionContainer") {
        return *p
    }
    return p.OtherItemMap[name]
}

func (p *SessionContainer)SetItem(item interface{}) {
    va := reflect.ValueOf(item);
    p.OtherItemMap[va.Type().String()] = item
}
func (p *SessionContainer)SetItemAlias(name string,item interface{}) {
    p.OtherItemMap[name] = item
}

func (p *SessionContainer)GetFuncParams(fun reflect.Value) (data []reflect.Value, err error) {
    var params []reflect.Value = nil;

    rf := fun.Type().String();

    ps := strings.Split(rf, "(")[1];
    ps = strings.Split(ps, ")")[0];

    //如果有参数需要注入
    if len(ps) != 0 {
        ps = strings.Replace(ps, " ", "", -1)
        paras := strings.Split(ps, ",");

        params = make([]reflect.Value, len(paras))

        for index, cla := range paras {
            p := p.GetItemByClassName(cla);
            if (p == nil) {
                return nil, errors.New("container not has '" +
                    cla + "' item , please use container.SetItem() to set item")
            }
            params[index] = reflect.ValueOf(p)
        }
    }

    return params, nil;
}


