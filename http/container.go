package http

import (
	"errors"
	"github.com/bysir-zl/bygo/cache"
	"github.com/bysir-zl/bygo/config"
	"github.com/bysir-zl/bygo/db"
	"lib.com/deepzz0/go-com/log"
	"reflect"
	"strings"
)

// 存储用于依赖注入的容器
type Container struct {
	DbFactory    db.DbFactory
	Cache        cache.CacheInterface
	Config       config.Config
	Logger       func(*Context, Logs)

	otherItemMap map[string]interface{}
}

func (p *Container) GetItemByClassName(name string) interface{} {
	return p.otherItemMap[name]
}

func (p *Container) GetItemByClass(item interface{}) interface{} {
	va := reflect.ValueOf(item)
	return p.otherItemMap[va.Type().String()]
}

// 向容器中添加一个Item
func (p *Container) SetItem(item interface{}) {
	va := reflect.ValueOf(item)
	p.otherItemMap[va.Type().String()] = item
}

//从容器中获取参数类型
func (s *Container) GetFuncParams(fun reflect.Value) (data []reflect.Value, err error) {
	var params []reflect.Value = nil

	rf := fun.Type().String()

	ps := strings.Split(rf, "(")[1]
	ps = strings.Split(ps, ")")[0]

	//如果有参数需要注入
	if len(ps) != 0 {
		ps = strings.Replace(ps, " ", "", -1)
		paras := strings.Split(ps, ",")

		params = make([]reflect.Value, len(paras))

		for index, cla := range paras {
			p := s.GetItemByClassName(cla)
			if p == nil {
				return nil, errors.New("container not has '" +
					cla + "' item , please use container.SetItem() to set item")
			}
			params[index] = reflect.ValueOf(p)
		}
	}

	return params, nil
}

func NewContainer() Container {
	return Container{
		otherItemMap: make(map[string]interface{}),
	}
}


///////////

// 一个请求的上下文
type Context struct {
	Response     *Response
	Request      *Request

	DbFactory    db.DbFactory
	Cache        cache.CacheInterface
	Config       config.Config
	Logger       func(*Context, Logs)

	otherItemMap map[string]interface{}
}

func (p *Context) GetItem(item interface{}) {
	va := reflect.TypeOf(item).Elem()
	key := va.PkgPath() + "#" + va.String()
	x := p.GetItemByAlias(key)
	v := reflect.ValueOf(item).Elem()

	if x != nil && v.CanSet() {
		v.Set(reflect.ValueOf(x))
	} else {
		log.Warn("GetItem '" + va.String() + "' error ! u forget SetItem or item is not a pointer . ")
	}
}

func (p *Context) GetItemByClass(item interface{}) interface{} {
	va := reflect.TypeOf(item)
	key := va.PkgPath() + "#" + va.String()

	return p.GetItemByAlias(key)
}

func (p *Context) GetItemByAlias(name string) interface{} {
	return p.otherItemMap[name]
}

func (p *Context) SetItem(item interface{}) {
	va := reflect.TypeOf(item)
	key := va.PkgPath() + "#" + va.String()

	p.otherItemMap[key] = item
}
func (p *Context) SetItemByAlias(name string, item interface{}) {
	p.otherItemMap[name] = item
}

func (p *Context) Resp(data ResponseData) {
	p.Response.Data = data
}

func (p *Context) SetBase(a *Context) {
	*p = *a
}

// 用于记录
func (p *Context) Log(code int, msg string) {
	if p.Logger != nil {
		p.Logger(p, Logs{Code:code, Message:msg})
	} else {
		log.Debugf("Code : %d , Msg : %s", code, msg)
	}
}

func NewContext() Context {
	return Context{
		otherItemMap: make(map[string]interface{}),
	}
}
