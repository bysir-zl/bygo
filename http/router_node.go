package http

import (
	"github.com/bysir-zl/bygo/bean"
	"reflect"
	"strings"

	"runtime"
	"lib.com/deepzz0/go-com/log"
)

var _ = log.Blue

type RouterNode struct {
	path           string       //当前path
	handlerType    string       //当前处理程序的类型
	handler        handler      //当前处理程序
	middlewareList []Middleware //当前的Middleware列表
	childrenList   []RouterNode //下一级
}
type handler struct {
	item           interface{}
	controllerFunc map[string]func(*Context) // 保存Controller的func，因为每次请求都反射会消耗性能
	handlerName    string
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
	routerNode.childrenList = []RouterNode{}
	routerNode.middlewareList = []Middleware{}
	routerNode.handlerType = "Group"

	p.childrenList = append(p.childrenList, routerNode)

	call(&routerNode)
	return &routerNode
}

//向当前节点添加中间件
func (p *RouterNode) Middleware(middleware Middleware) *RouterNode {
	p.middlewareList = append(p.middlewareList, middleware)

	return p
}

//在当前节点添加一个处理控制器的子节点
func (p *RouterNode) Controller(path string, controller interface{}) *RouterNode {
	path = formatPath(path)
	//新建一个子node
	routerNode := RouterNode{}

	routerNode.path = path
	routerNode.handlerType = "Controller"
	routerNode.middlewareList = []Middleware{}
	routerNode.handler.item = controller
	routerNode.handler.controllerFunc = map[string]func(*Context){}
	routerNode.handler.handlerName = reflect.ValueOf(controller).Type().String()

	//controller.()
	stru := reflect.ValueOf(controller)
	typ := stru.Type()

	// 取出所有的方法, 检查签名, 若签名正确就保存到map里
	for i := stru.NumMethod() - 1; i >= 0; i-- {
		fun := stru.Method(i)
		ifun, ok := fun.Interface().(func(*Context))

		if ok {
			name := typ.Method(i).Name
			// 省略 OMIT
			name = strings.TrimPrefix(name, "OMIT")
			routerNode.handler.controllerFunc[name] = ifun
		}
	}

	p.childrenList = append(p.childrenList, routerNode)
	return &routerNode
}

//在当前节点添加一个处理函数的子节点
func (p *RouterNode) Fun(path string, fun func(*Context)) *RouterNode {
	path = formatPath(path)

	//新建一个子node
	routerNode := RouterNode{}

	routerNode.handler.item = fun
	routerNode.path = path
	routerNode.handlerType = "Func"
	routerNode.middlewareList = []Middleware{}
	funcInfo := runtime.FuncForPC(reflect.ValueOf(fun).Pointer()).Name()

	funcInfo = strings.Replace(funcInfo, "-fm", "", -1)
	funcInfos := strings.Split(funcInfo, ".")
	funcInfo = strings.Join(funcInfos[len(funcInfos) - 2:], ".")
	funcInfo = strings.Replace(funcInfo, ")", "", -1)
	funcInfo = strings.Replace(funcInfo, "(", "", -1)

	routerNode.handler.handlerName = funcInfo

	p.childrenList = append(p.childrenList, routerNode)
	return &routerNode
}

func (node *RouterNode) run(context *Context, otherUrl string) {

	request := context.Request
	response := context.Response

	if node.handlerType == "Controller" {
		method := "Index"

		// 解析方法与路由参数
		if otherUrl != "" {
			urlParamsList := strings.Split(otherUrl, "/")

			if len(urlParamsList) > 0 {
				method = urlParamsList[0]
				// 大写第一个字母
				method = strings.ToUpper(string(method[0])) + string(method[1:])

				if node.handler.controllerFunc[method] == nil && node.handler.controllerFunc["Index"] != nil {
					method = "Index"
					request.Router.Params = urlParamsList
				} else {
					if len(urlParamsList) > 1 {
						request.Router.Params = urlParamsList[1:]
					}
				}
			}
		}
		request.Router.Handler = "controller|" + node.handler.handlerName + "." + method
		request.Router.Method = method

		// 从controller中读取一个方法
		fun := node.handler.controllerFunc[method]
		//没找到类方法,url不正确
		if fun == nil {
			msg := "the method '" + method +
				"' is undefined in controller '" + node.handler.handlerName + "'!"
			response.Data = NewResponseJson(404, bean.ApiData{Code: 404, Msg: msg})
			return
		}
		fun(context)
		return
	} else if node.handlerType == "Func" {
		if otherUrl != "" {
			request.Router.Params = strings.Split(otherUrl, "/")
		}

		fun := node.handler.item.(func(*Context))
		request.Router.Handler = "func|" + node.handler.handlerName
		fun(context)
		return
	} else {
		//没有配置路由
		response.Data = NewResponseJson(404, bean.ApiData{Code: 404, Msg: "u are forget set route? but welcome use bygo . :D"})
		return
	}

	return
}
