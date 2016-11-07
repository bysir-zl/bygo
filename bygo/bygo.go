package bygo

import (
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
	AppContainer byhttp.Container

	Router       byhttp.Router
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
	context.Logger = p.AppContainer.Logger

	p.Router.Start(r.URL.String(), &context)

	if context.Response.Result == nil {
		context.Response.Result.Apply(r, w)
	} else {
		w.Header().Set("Content-Type", context.Response.Data.Type)
		w.WriteHeader(context.Response.Data.Code)
		io.WriteString(w, context.Response.Data.Body)
	}
}

func (p *ApiHandler) ConfigRouter(root string, fun func(*byhttp.RouterNode)) {
	p.Router.Init(root,fun)
}

func (p *ApiHandler) ConfigLogger(fun func(*byhttp.Context, byhttp.Logs)) {
	p.AppContainer.Logger = fun
}

func Config(files ...string) {
	for _, file := range files {

		if _, er := os.Stat(file); er == nil || os.IsExist(er) {
			c, err := config.LoadConfigFromFile(file)
			if err != nil {
				log.Warn(err)
			}
			bConfig = c
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
	apiHandle = &ApiHandler{
		AppContainer: byhttp.NewContainer(),
		Router:       byhttp.NewRouter(),
	}
	return
}
