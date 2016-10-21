package http

type RouterNode struct {
    Path           string        //当前path
    HandlerType    string        //当前处理程序的类型
    Handler        interface{}   //当前处理程序
    MiddlewareList *[]Middleware //当前的Middleware列表
    ChildrenList   *[]RouterNode //下一级
}

func formatPath(path string) string {
    if (path==""){
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