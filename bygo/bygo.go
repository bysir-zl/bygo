package bygo

import (
    "net/http"
    byhttp "github.com/bysir-zl/bygo/http"
    "io"
    "github.com/bysir-zl/bygo/cache"
    "github.com/bysir-zl/bygo/db"
    "github.com/bysir-zl/bygo/bean"
    "config"
)

type ApiHandler struct {
    AppContainer    byhttp.Container

    Router          byhttp.Router
    ExceptionHandle func(byhttp.SessionContainer, byhttp.Exceptions) byhttp.ResponseData
    Debug           bool
}

func NewApiHandler() (apiHandle *ApiHandler) {
    //App 容器
    appContainer := byhttp.Container{
        OtherItemMap:make(map[string]interface{}),
    };

    node := byhttp.RouterNode{}
    node.ChildrenList = &[]byhttp.RouterNode{}
    node.MiddlewareList = &[]byhttp.Middleware{}
    node.HandlerType = "Base"

    apiHandle = &ApiHandler{
        AppContainer:appContainer,
        Router:byhttp.Router{RouterNode: node},
    }
    return
}

func (p *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    response := byhttp.Response{ResponseWrite: w};

    request := byhttp.Request{Request:r};
    request.Init();

    sessionContainer := byhttp.SessionContainer{
        OtherItemMap:make(map[string]interface{}),
        Request:&request,
        Response:&response,
    };

    //将app容器添加到session容器
    sessionContainer.Cache = p.AppContainer.Cache
    sessionContainer.DbFactory = p.AppContainer.DbFactory
    for k, v := range p.AppContainer.OtherItemMap {
        sessionContainer.OtherItemMap[k] = v
    }

    responseData := p.Router.Start(r.URL.String(), sessionContainer);

    //错误处理
    if (responseData.Exception.Code != 0) {
        if p.ExceptionHandle != nil {
            responseData = p.ExceptionHandle(sessionContainer, responseData.Exception)
        } else {
            responseData = byhttp.NewRespDataJson(
                responseData.Code,
                bean.ApiData{
                    Code:responseData.Exception.Code,
                    Msg:responseData.Exception.Message,
                })
        }
    }

    w.Header().Set("Content-Type", responseData.Type);
    w.WriteHeader(responseData.Code)
    io.WriteString(w, responseData.Body)
}

func (p *ApiHandler) ConfigRouter(root string, fun func(*byhttp.RouterNode)) {
    p.Router.RouterNode.Root(root)
    fun(&p.Router.RouterNode)
}

func (p *ApiHandler) ConfigExceptHandler(fun func(byhttp.SessionContainer, byhttp.Exceptions) byhttp.ResponseData) {
    p.ExceptionHandle = fun
}

func (p *ApiHandler) Init() {
    c := cache.NewCache(config.CacheDriver)
    dbFactory := db.NewDbFactory(config.DbConfigs)

    p.AppContainer.Cache = c
    p.AppContainer.DbFactory = dbFactory
}