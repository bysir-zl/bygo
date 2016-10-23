package http

import (
    "reflect"
    "errors"
    "log"
)

type RouterNode struct {
    Path           string                   //当前path
    HandlerType    string                   //当前处理程序的类型
    Handler        interface{}              //当前处理程序
    MiddlewareList *[]Middleware            //当前的Middleware列表
    ChildrenList   *[]RouterNode            //下一级

    ControllerFunc map[string]reflect.Value // 保存Controller的func，因为每次请求都反射会消耗性能
}

func formatPath(path string) string {
    if (path == "") {
        return ""
    }
    if path[0] != '/' {
        path = "/" + path;
    }
    if path[len(path) - 1] == '/' {
        path = path[:len(path) - 1];
    }
    return path
}

func (p *RouterNode) Root(path string) {
    path = formatPath(path)
    p.Path = path
}

// 当前节点添加一个子节点,并传递子节点调用方法
// 用于嵌套路由
func (p *RouterNode) Group(path string, call func(*RouterNode)) *RouterNode {
    path = formatPath(path)

    //新建一个子node
    routerNode := RouterNode{};
    routerNode.Path = path
    routerNode.ChildrenList = &[]RouterNode{}
    routerNode.MiddlewareList = &[]Middleware{}
    routerNode.HandlerType = "Group"
    *p.ChildrenList = append(*p.ChildrenList, routerNode);

    call(&routerNode)
    return &routerNode
}

//向当前节点添加中间件
func (p *RouterNode) Middleware(middleware Middleware) *RouterNode {
    *p.MiddlewareList = append(*p.MiddlewareList, middleware);

    return p
}

//在当前节点添加一个处理控制器的子节点
func (p *RouterNode) Controller(path string, controller interface{}) *RouterNode {
    path = formatPath(path)


    //新建一个子node
    routerNode := RouterNode{};

    routerNode.Path = path
    routerNode.Handler = controller
    routerNode.HandlerType = "Controller"
    routerNode.MiddlewareList = &[]Middleware{}
    routerNode.ControllerFunc = map[string]reflect.Value{}

    stru := reflect.ValueOf(controller).Elem();
    for i := stru.NumMethod() - 1; i > 0; i-- {
        fun := stru.Method(i)
        routerNode.ControllerFunc[fun.Type().Name()] = fun
    }

    //log.Print(controller)

    log.Print(routerNode.ControllerFunc)

    *p.ChildrenList = append(*p.ChildrenList, routerNode);
    return &routerNode
}

//在当前节点添加一个处理函数的子节点
func (p *RouterNode) Get(path string, fun interface{}) *RouterNode {
    path = formatPath(path)

    //新建一个子node
    routerNode := RouterNode{};

    routerNode.Path = path
    routerNode.Handler = fun
    routerNode.HandlerType = "Func"
    routerNode.MiddlewareList = &[]Middleware{}

    *p.ChildrenList = append(*p.ChildrenList, routerNode);
    return &routerNode
}

//2016/8/27
//在当前节点添加一个处理Model的子节点
func (p *RouterNode) Model(path string, model RouterModelInterface) *RouterNode {
    path = formatPath(path)

    //新建一个子node
    routerNode := RouterNode{};

    routerNode.Path = path
    routerNode.Handler = model
    routerNode.HandlerType = "Model"
    routerNode.MiddlewareList = &[]Middleware{}

    *p.ChildrenList = append(*p.ChildrenList, routerNode);
    return &routerNode
}

func (node *RouterNode) run(sessionContainer SessionContainer, method string) (handlerName string, response ResponseData) {
    var fun reflect.Value = reflect.Value{};

    if node.HandlerType == "Model" {
        modelHandler := RouterModelHandler{
            model:node.Handler.(RouterModelInterface),
            method:method,
        };
        //运行Model的Handle
        response = modelHandler.Handle(sessionContainer);
        return

    } else if node.HandlerType == "Controller" {
        // 从controller中读取一个方法
        // 这里可以优化为 程序运行时就将fun读到map中
        handlerName = reflect.ValueOf(node.Handler).Type().Name() + "@" + method
        fun = node.ControllerFunc[method]

        //没找到类方法,url不正确
        var zero reflect.Value
        if fun ==   zero {
            response = NewRespDataError(404, errors.New("the method '" + method + "' is undefined " +
                "in controller '" + reflect.TypeOf(node.Handler).String() + "'!"))
            return
        }
    } else if node.HandlerType == "Func" {
        fun = reflect.ValueOf(node.Handler)
        handlerName = "func"
    } else {
        //没有配置路由
        response = NewRespDataError(500, errors.New("u are forget set route? bug welcome use bygo . :D"))
        return
    }

    //从容器中获取参数
    params, err := sessionContainer.GetFuncParams(fun);
    if (err != nil) {
        response = NewRespDataError(500, err)
        return
    }

    response = (fun.Call(params)[0]).Interface().(ResponseData)
    return
}
