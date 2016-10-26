package http

import (
    "strings"
    "log"
    "runtime/debug"
    "reflect"
)

type Router struct {
    RouterNode RouterNode
    RouterPath map[string][]RouterNode // 根据router设置解析出来的节点列表
}

// 根据url匹配路由,
// url:当前路径
// node:当前路由节点
// nodeList:当前递归路径
// currNodeList:匹配到的node路径
// deep:当前递归深度
// currDeep:匹配到的深度
func (p *Router)Handler(allUrl string) (matchedUrl string, matchedNodeList []RouterNode) {
    if (allUrl == "") {
        return
    }

    var pathMaxLen = 0
    for path, nodes := range p.RouterPath {
        if strings.Index(allUrl, path) >= 0 {
            pathLen := len(path)
            if pathLen > pathMaxLen {
                pathMaxLen = pathLen
                matchedNodeList = nodes
                matchedUrl = path
            }
        }
    }

    return

    //  old slow code
    //index := strings.Index(url, node.Path)
    //
    //// 以 /xx开头,并且(/结尾 或者 空结尾)
    //// 或者 path=="" 则默认匹配成功,继续匹配下一级
    //if node.Path == "" || (index == 0 && (len(url) == len(node.Path) + 1 || url[len(node.Path)] == '/')) {
    //    *nodeList = append(*nodeList, *node)
    //
    //    if node.Path != "" || node.HandlerType != "Group" {
    //        // 如果不是空组(空组可能是为了分组设置中间件) , 深度才加一
    //        // 深度会影响到最深路径的判断
    //        deep = deep + 1
    //    }
    //
    //    //log.Println(url, "'"+node.Path+"'", node.HandlerType)
    //
    //    // 处理器全匹配直接返回不再匹配其他
    //    if (node.HandlerType == "Func" || node.HandlerType == "Model" || node.HandlerType == "Controller") &&
    //        len(url) == len(node.Path) + 1 {
    //        *matchedNodeList = *nodeList
    //        return false
    //    }
    //
    //    // 保存匹配到的最深的路径
    //    if (deep > matchedDeep) {
    //        matchedDeep = deep
    //        *matchedNodeList = *nodeList
    //    }
    //
    //    if (node.ChildrenList != nil) {
    //        u := url[len(node.Path):];
    //        for _, children := range *node.ChildrenList {
    //            if !p.Handler(allUrl, u, &children, nodeList, matchedNodeList, deep, matchedDeep) {
    //                break
    //            }
    //        }
    //    }
    //}
    //return true
}

func (p *Router)ParseToPath(matchedUrl string, node *RouterNode, nodeList *[]RouterNode) {

    matchedUrl = matchedUrl + node.Path

    *nodeList = append(*nodeList, *node)

    if node.ChildrenList != nil {
        for _, children := range *node.ChildrenList {
            p.ParseToPath(matchedUrl, &children, nodeList)
        }
    } else {
        p.RouterPath[matchedUrl] = *nodeList
    }
}

func (p *Router)Init(fun func(node *RouterNode)) {
    fun(&p.RouterNode)

    nodeList := []RouterNode{}
    p.ParseToPath("", &p.RouterNode, &nodeList);
}

func (p *Router)Start(url string, sessionContainer Context) (ResponseData) {

    defer func() {
        if err := recover(); err != nil {
            log.Println("-----------ERROR---------------")
            log.Println(err)
            debug.PrintStack()
        }
    }()

    request := sessionContainer.Request;
    response := sessionContainer.Response;

    baseUrl := strings.Split(strings.Split(url, "?")[0], "#")[0];
    // 加上/以匹配"/"根
    if (baseUrl[len(baseUrl) - 1] != '/') {
        baseUrl = baseUrl + "/"
    }

    urs := strings.Split(url, "#");
    urlHash := ""
    if (len(urs) > 1) {
        urlHash = urs[1]
    }

    matchedUrl, currNodeList := p.Handler(baseUrl);

    var node RouterNode;
    // 没有匹配到东西
    if len(currNodeList) == 0 {
        node = p.RouterNode
    } else {
        node = currNodeList[len(currNodeList) - 1]
    }

    otherParam := string(baseUrl[len(matchedUrl):])

    //去掉前后多余的/
    otherParam = strings.TrimLeft(otherParam, "/")
    otherParam = strings.TrimRight(otherParam, "/")

    //log.Println(baseUrl, matchedUrl, otherParam)

    method := "Index"; //默认进入Index方法

    if (otherParam != "") {
        urlParamsList := strings.Split(otherParam, "/");

        if node.HandlerType == "Controller" || node.HandlerType == "Model" {
            // 有参
            if (len(urlParamsList) > 0) {
                method = urlParamsList[0]
                // 大写第一个字母
                method = strings.ToUpper(string(method[0])) + string(method[1:])
            }
            // 除了第一个作为方法名,还有多余的参
            if (len(urlParamsList) > 1) {
                urlParams := urlParamsList[1:]
                request.Router.Params = urlParams;
            }
        } else if (node.HandlerType == "Func") {
            request.Router.Params = urlParamsList;
        }
    }

    request.Router.Url = url
    request.Router.Hash = urlHash

    // 运行中间件
    stop := false
    for _, item := range currNodeList {
        if (item.MiddlewareList != nil) {
            for _, item := range *item.MiddlewareList {
                needStop, data := item.HandlerBefore(sessionContainer)
                if (needStop) {
                    stop = true
                    response.ResponseData = data
                    request.Router.Handler = reflect.TypeOf(item).Name()
                }
            }
        }
    }

    if !stop {
        // 运行某个node
        request.Router.Handler, response.ResponseData = node.run(sessionContainer, method);
    }

    // 倒着运行中间件
    for i := len(currNodeList) - 1; i >= 0; i = i - 1 {
        item := currNodeList[i]
        if (item.MiddlewareList != nil) {
            for _, item := range *item.MiddlewareList {
                needStop, data := item.HandlerAfter(sessionContainer)
                if (needStop) {
                    return data;
                }
            }
        }
    }

    return response.ResponseData
}

func NewRouter() Router {
    node := RouterNode{}
    node.ChildrenList = &[]RouterNode{}
    node.MiddlewareList = &[]Middleware{}
    node.HandlerType = "Base"
    return Router{
        RouterNode: node,
        RouterPath:map[string][]RouterNode{},
    }
}