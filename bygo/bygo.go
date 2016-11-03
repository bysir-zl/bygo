package bygo

import (
	"github.com/bysir-zl/bygo/bean"
	"github.com/bysir-zl/bygo/cache"
	"github.com/bysir-zl/bygo/config"
	"github.com/bysir-zl/bygo/db"
	byhttp "github.com/bysir-zl/bygo/http"
	"io"
	"lib.com/deepzz0/go-com/log"
	"net/http"
	"os"
)

var _ = log.Blue
var bConfig = config.Config{}

type ApiHandler struct {
	AppContainer    byhttp.Container

	Router          byhttp.Router
	ExceptionHandle func(*byhttp.Context)
}

func (p *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := byhttp.Response{ResponseWrite: w}

	request := byhttp.Request{Request: r}
	request.Init()

	context := byhttp.NewContext()
	context.Request = &request
	context.Response = &response

	//将app容器添加到session容器
	context.Cache = p.AppContainer.Cache
	context.DbFactory = p.AppContainer.DbFactory
	context.Config = p.AppContainer.Config

	for k, v := range p.AppContainer.OtherItemMap {
		context.SetItemByAlias(k, v)
	}

	p.Router.Start(r.URL.String(), &context)

	// 错误处理
	// todo 这里没有运行中间件,是一个bug
	if context.Response.Data.Exception.Code != 0 {
		if p.ExceptionHandle != nil {
			p.ExceptionHandle(&context)
		} else {
			context.Response.Data = byhttp.NewRespDataJson(
				context.Response.Data.Code,
				bean.ApiData{
					Code: context.Response.Data.Exception.Code,
					Msg:  context.Response.Data.Exception.Message,
				})
		}
	}

	w.Header().Set("Content-Type", context.Response.Data.Type)
	w.WriteHeader(context.Response.Data.Code)
	io.WriteString(w, context.Response.Data.Body)
}

func (p *ApiHandler) ConfigRouter(root string, fun func(*byhttp.RouterNode)) {
	p.Router.RouterNode.Root(root)
	p.Router.Init(fun)
}

func (p *ApiHandler) ConfigExceptHandler(fun func(*byhttp.Context)) {
	p.ExceptionHandle = fun
}
func Config(files ...string) {
	for _, file := range files {

		if _, er := os.Stat(file); er == nil || os.IsExist(er) {
			config, err := config.LoadConfigFromFile(file)
			if err != nil {
				log.Warn(err)
			}
			bConfig = config
			return
		}
	}

}

func (p *ApiHandler) Init() {
	if bConfig.Evn == "" {
		log.Warn("you have not config the bygo , please use bygo.Config(filePath) to config it .")
	}

	c := cache.NewCache(bConfig)
	dbFactory := db.Init(bConfig.DbConfigs)

	p.AppContainer.Cache = c
	p.AppContainer.DbFactory = dbFactory
	p.AppContainer.Config = bConfig

	log.Info("apiHandler evn is : " + p.AppContainer.Config.Evn)
}

func NewApiHandler() (apiHandle *ApiHandler) {
	//App 容器
	appContainer := byhttp.Container{
		OtherItemMap: make(map[string]interface{}),
	}

	apiHandle = &ApiHandler{
		AppContainer: appContainer,
		Router:       byhttp.NewRouter(),
	}
	return
}
