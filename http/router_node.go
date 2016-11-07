package http

import (
	"github.com/bysir-zl/bygo/bean"
	"reflect"
	"strings"
)

type RouterNode struct {
	path           string            //当前path
	handlerType    string            //当前处理程序的类型
	handler        interface{}       //当前处理程序
	middlewareList *[]Middleware     //当前的Middleware列表
	childrenList   *[]RouterNode     //下一级

	controllerFunc map[string]func() // 保存Controller的func，因为每次请求都反射会消耗性能
}

func formatPath(path string) string {
	if path == "" {
		return ""
	}
	if path[0] != '/' {
		path = "/" + path
	}
	if path[len(path) - 1] == '/' {
		path = path[:len(path) - 1]
	}
	return path
}

// 当前节点添加一个子节点,并传递子节点调用方法
// 用于嵌套路由
func (p *RouterNode) Group(path string, call func(*RouterNode)) *RouterNode {
	path = formatPath(path)

	//新建一个子node
	routerNode := RouterNode{}
	routerNode.path = path
	routerNode.childrenList = &[]RouterNode{}
	routerNode.middlewareList = &[]Middleware{}
	routerNode.handlerType = "Group"
	*p.childrenList = append(*p.childrenList, routerNode)

	call(&routerNode)
	return &routerNode
}

//向当前节点添加中间件
func (p *RouterNode) Middleware(middleware Middleware) *RouterNode {
	*p.middlewareList = append(*p.middlewareList, middleware)

	return p
}

//在当前节点添加一个处理控制器的子节点
func (p *RouterNode) Controller(path string, controller ControllerInterface) *RouterNode {
	path = formatPath(path)
	//新建一个子node
	routerNode := RouterNode{}

	routerNode.path = path
	routerNode.handler = controller
	routerNode.handlerType = "Controller"
	routerNode.middlewareList = &[]Middleware{}
	routerNode.controllerFunc = map[string]func(){}
	//controller.()
	stru := reflect.ValueOf(controller)
	typ := stru.Type()

	// 取出所有的方法, 检查签名, 若签名正确就保存到map里
	for i := stru.NumMethod() - 1; i >= 0; i-- {
		fun := stru.Method(i)
		ifun, ok := fun.Interface().(func())

		if ok {
			routerNode.controllerFunc[typ.Method(i).Name] = ifun
		}
	}

	*p.childrenList = append(*p.childrenList, routerNode)
	return &routerNode
}

//在当前节点添加一个处理函数的子节点
func (p *RouterNode) Fun(path string, fun func(*Context)) *RouterNode {
	path = formatPath(path)

	//新建一个子node
	routerNode := RouterNode{}

	routerNode.path = path
	routerNode.handler = fun
	routerNode.handlerType = "Func"
	routerNode.middlewareList = &[]Middleware{}

	*p.childrenList = append(*p.childrenList, routerNode)
	return &routerNode
}

func (node *RouterNode) run(context *Context, otherUrl string) {

	request := context.Request
	response := context.Response

	if node.handlerType == "Controller" {
		node.handler.(ControllerInterface).SetBase(context)
		method := "Index"

		// 解析方法与路由参数
		if otherUrl != "" {
			urlParamsList := strings.Split(otherUrl, "/")

			if len(urlParamsList) > 0 {
				method = urlParamsList[0]
				// 大写第一个字母
				method = strings.ToUpper(string(method[0])) + string(method[1:])

				if node.controllerFunc[method] == nil && node.controllerFunc["Index"] != nil {
					method = "Index"
					request.Router.Params = urlParamsList
				} else {
					if len(urlParamsList) > 1 {
						request.Router.Params = urlParamsList[1:]
					}
				}
			}
		}
		request.Router.Handler = reflect.ValueOf(node.handler).Type().String() + "@" + method
		// 从controller中读取一个方法
		fun := node.controllerFunc[method]
		//没找到类方法,url不正确
		if fun == nil {
			msg := "the method '" + method +
				"' is undefined in controller '" + reflect.TypeOf(node.handler).String() + "'!"
			response.Data = NewRespDataJson(404, bean.ApiData{Code: 404, Msg: msg})
			return
		}
		fun()
		return
	} else if node.handlerType == "Func" {
		if otherUrl != "" {
			request.Router.Params = strings.Split(otherUrl, "/")
		}

		fun := node.handler.(func(*Context))
		request.Router.Handler = "func"
		fun(context)
		return
	} else {
		//没有配置路由
		response.Data = NewRespDataJson(404, bean.ApiData{Code: 404, Msg: "u are forget set route? but welcome use bygo . :D"})
		return
	}

	return
}
