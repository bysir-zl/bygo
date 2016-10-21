package http

import (
    "reflect"
    "strings"
    "log"
    "runtime/debug"
    "errors"
)

type Router struct {
    RouterNode RouterNode
}

// 根据url匹配路由,
// url:全路径
// node:路由节点
// nodeList:当前递归路径
// currNodeList:匹配到的node路径
// deep:当前递归深度
// currDeep:匹配到的深度
func Handler(url string, node *RouterNode, nodeList *[]RouterNode, currNodeList *[]RouterNode, deep int, currDeep int) (hasNext bool) {
    if (url == "") {
        return false
    }

    index := strings.Index(url, node.Path)

    // 以 /xx开头,并且(/结尾 或者 空结尾)
    // 或者 path=="" 则默认匹配成功,继续匹配下一级
    if node.Path == "" || (index == 0 && (len(url) == len(node.Path) + 1 || url[len(node.Path)] == '/')) {
        *nodeList = append(*nodeList, *node)

        if node.Path != "" || node.HandlerType != "Group" {
            // 如果不是空组(空组可能是为了分组设置中间件) , 深度才加一
            // 深度会影响到最深路径的判断
            deep = deep + 1
        }

        //log.Println(url, "'"+node.Path+"'", node.HandlerType)

        // 处理器全匹配直接返回不再匹配其他
        if (node.HandlerType == "Func" || node.HandlerType == "Model" || node.HandlerType == "Controller") &&
            len(url) == len(node.Path) + 1 {
            *currNodeList = *nodeList
            return false
        }

        // 保存匹配到的最深的路径
        if (deep > currDeep) {
            currDeep = deep
            *currNodeList = *nodeList
        }

        if (node.ChildrenList != nil) {
            u := url[len(node.Path):];
            for _, children := range *node.ChildrenList {
                if !Handler(u, &children, nodeList, currNodeList, deep, currDeep) {
                    break
                }
            }
        }
    }
    return true
}

func (p *Router)Start(url string, sessionContainer SessionContainer) (ResponseData) {

    defer func() {
        if err := recover(); err != nil {
            log.Println("-----------ERROR---------------")
            log.Println(err)
            debug.PrintStack()
        }
    }()

    request := sessionContainer.Request;

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

    var nodeList []RouterNode = []RouterNode{}
    var currNodeList []RouterNode = []RouterNode{}

    Handler(baseUrl, &p.RouterNode, &nodeList, &currNodeList, 0, 0);

    //获取最后一个Handler,就是成功匹配到的Handler
    var node RouterNode;
    matchedUrl := "";
    var middlewareList []Middleware = []Middleware{};
    for _, item := range currNodeList {
        node = item;
        matchedUrl = matchedUrl + item.Path;
        // 将当前node中的中间件一次加载到要运行的middlewareList中
        if (item.MiddlewareList != nil) {
            for _, middlewareItem := range *item.MiddlewareList {
                middlewareList = append(middlewareList, middlewareItem);
            }
        }
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
            //有参
            if (len(urlParamsList) > 0) {
                method = urlParamsList[0]
                //大写第一个字母
                method = strings.ToUpper(string(method[0])) + string(method[1:])
            }
            //除了第一个作为方法名,还有多余的参
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
    request.Router.Handler = reflect.TypeOf(node.Handler).Name() + "@" + method

    //运行中间件
    for _, item := range middlewareList {
        needStop, data := item.HandlerBefore(sessionContainer)
        if (needStop) {
            return data;
        }
    }

    // 处理运行某个node

    var fun reflect.Value = reflect.Value{};
    if node.HandlerType == "Model" {
        modelHandler := RouterModelHandler{
            model:node.Handler.(RouterModelInterface),
            method:method,
        };
        //运行Model的Handle
        return modelHandler.Handle(sessionContainer);

    } else if node.HandlerType == "Controller" {
        // 从controller中读取一个方法
        fv := reflect.ValueOf(node.Handler).Elem();
        me := fv.MethodByName(method)

        //没找到类方法,url不正确
        if !me.IsValid() {
            return NewRespDataError(404, errors.New("the method '" + method + "' is undefined " +
                "in controller '" + reflect.TypeOf(node.Handler).String() + "'!"))
        }
        fun = me;
    } else if node.HandlerType == "Func" {
        fun = reflect.ValueOf(node.Handler)
    } else {
        //没有配置路由
        return ResponseData{Code:200, Body:"<h1>Welcome Use Bygo</h1>"}
    }

    //从容器中获取参数
    params, err := sessionContainer.GetFuncParams(fun);
    if (err != nil) {
        return NewRespDataError(500, err)
    }

    //执行方法
    sessionContainer.Response.ResponseData = (fun.Call(params)[0]).Interface().(ResponseData)

    //倒着运行中间件
    l := len(middlewareList)

    for i := l - 1; i >= 0; i = i - 1 {
        item := middlewareList[i]
        needStop, data := item.HandlerAfter(sessionContainer)
        if (needStop) {
            return data;
        }
    }

    return sessionContainer.Response.ResponseData
}

