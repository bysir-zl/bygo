package bygo

import (
    "net/http"
    byhttp "github.com/bysir-zl/bygo/http"
    "io"
    "github.com/bysir-zl/bygo/cache"
    "github.com/bysir-zl/bygo/db"
    "github.com/bysir-zl/bygo/bean"
    "os"
    "github.com/deepzz0/go-com/log"
)

var _ = log.Blue

type ApiHandler struct {
    AppContainer    byhttp.Container

    Router          byhttp.Router
    ExceptionHandle func(byhttp.SessionContainer, byhttp.Exceptions) byhttp.ResponseData
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
    sessionContainer.Config = p.AppContainer.Config

    for k, v := range p.AppContainer.OtherItemMap {
        sessionContainer.OtherItemMap[k] = v
    }

    responseData := p.Router.Start(r.URL.String(), sessionContainer);

    // 错误处理
    // todo 这里没有运行中间件,是一个bug
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
    p.Router.Init(fun)
}

func (p *ApiHandler) ConfigExceptHandler(fun func(byhttp.SessionContainer, byhttp.Exceptions) byhttp.ResponseData) {
    p.ExceptionHandle = fun
}
func (p *ApiHandler) Config(files ...string) {
    for _, file := range files {

        if _, er := os.Stat(file); er == nil || os.IsExist(er) {
            p.AppContainer.Config.LoadFromFile(file)
            return
        }
    }
}

func (p *ApiHandler) Init() {
    c := cache.NewCache(p.AppContainer.Config.CacheDrive)
    dbFactory := db.NewDbFactory(p.AppContainer.Config.DbConfigs)

    p.AppContainer.Cache = c
    p.AppContainer.DbFactory = dbFactory

    db.BFactory = dbFactory

    log.Info("apiHandler evn is : " + p.AppContainer.Config.Evn)
}

func NewApiHandler() (apiHandle *ApiHandler) {
    //App 容器
    appContainer := byhttp.Container{
        OtherItemMap:make(map[string]interface{}),
    };

    apiHandle = &ApiHandler{
        AppContainer:appContainer,
        Router:byhttp.NewRouter(),
    }
    return
}