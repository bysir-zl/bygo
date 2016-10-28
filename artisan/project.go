package artisan

import (
    "strings"
    "os"
    "log"
)

func checkFileIsExist(filename string) (bool) {
    var exist = true;
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        exist = false;
    }
    return exist;
}

func CreateProject(name string) {

    pathApp := "controller,exception,middleware,model,router,tool"

    for _, path := range strings.Split(pathApp, ",") {
        path = name + "/src/app/" + path;

        os.MkdirAll(path, os.ModePerm)
    }
    os.MkdirAll(name + "/src/config", os.ModePerm)

    files := map[string]string{}

    // Index
    files["index.go"] = "package main\n\nimport (\n    \"app/exception\"\n    \"app/router\"\n    \"github.com/bysir-zl/bygo/bygo\"\n    \"log\"\n    \"net/http\"\n)\n\nfunc main() {\n\n    apiHandle := bygo.NewApiHandler()\n\n    apiHandle.ConfigRouter(\"api\", router.Init)\n    apiHandle.ConfigExceptHandler(exception.Handler)\n    apiHandle.Config(\"./config/debug.json\")\n    apiHandle.Init()\n\n    http.Handle(\"/api/\", apiHandle)\n    http.Handle(\"/\", http.FileServer(http.Dir(\"./dist\")))\n\n    log.Println(\"server start success\")\n\n    err := http.ListenAndServe(\":81\", nil)\n\n    if err != nil {\n        log.Println(err)\n    }\n}\n"
    // router
    files["app/router/router.go"] = "package router\n\nimport (\n    \"app/middleware\"\n    \"app/controller\"\n    \"github.com/bysir-zl/bygo/http\"\n)\n\nfunc Init(node *http.RouterNode) {\n    node.Middleware(&middleware.HeaderMiddleware{}) // 为当前节点添加中间件\n\n    node.Get(\"/\", func(request *http.Request, p *http.Response) http.ResponseData {\n        return http.NewRespDataHtml(404, \"blank\")\n    })\n\n    node.Controller(\"index\", &controller.IndexController{})\n}\n"
    // IndexController
    files["app/controller/IndexController.go"] = "package controller\n\nimport (\n    \"fmt\"\n    \"strings\"\n    \"github.com/bysir-zl/bygo/http\"\n    \"strconv\"\n)\n\nvar count = 0\n\ntype IndexController struct{}\n\nfunc (p IndexController) Index(r *http.Request, s *http.Response) http.ResponseData {\n    count++\n\n    return http.NewRespDataHtml(200, \"welcome to use bygo!\" + \"<br>\" +\n        \"Url: \" + r.Router.Url + \"<br>\" +\n        \"Handler: \" + r.Router.Handler + \"<br>\" +\n        \"RouterParams : \" + strings.Join(r.Router.Params, \",\") + \"<br>\" +\n        \"Input : \" + fmt.Sprint(r.Input.All()) + \"<br>\" +\n        \"Header : \" + fmt.Sprint(r.Header) + \"<br>\" +\n        \"Count : \" + strconv.Itoa(count) + \"<br>\" +\n        \"\")\n}\n"

    // HeaderMiddleware
    files["app/middleware/HeaderMiddleware.go"] = "package middleware\n\nimport (\n    \"github.com/bysir-zl/bygo/http\"\n)\n\ntype HeaderMiddleware struct {\n\n}\n\nfunc (p HeaderMiddleware) HandlerBefore(s http.Context) (needStop bool, data http.ResponseData) {\n    s.Request.Input.Set(\"ext\", \"from middleware\")\n    return false, http.NewRespDataHtml(0, \"\")\n}\n\nfunc (p HeaderMiddleware) HandlerAfter(s http.Context) (needStop bool, data http.ResponseData) {\n\n    response := s.Response\n    response.ResponseData.Body = response.ResponseData.Body + \"<br><br> i am from middleware\"\n    response.AddHeader(\"Access-Control-Allow-Origin\", \"*\") // 添加上允许跨域\n    response.AddHeader(\"Access-Control-Allow-Headers\", \"X_TOKEN\") // 添加上允许的头,用来身份验证\n\n    return false, http.NewRespDataHtml(0, \"\")\n}\n"

    // HeaderMiddleware
    files["app/exception/Exception.go"] = "package exception\n\nimport (\n    \"github.com/bysir-zl/bygo/bean\"\n    \"github.com/bysir-zl/bygo/http\"\n)\n\n// 将报错的Exception处理成Response返回。在这里你可以判断e.Code统一处理错误,比如上报code==500的错误\nfunc Handler(c http.Context, e http.Exceptions) http.ResponseData {\n    return http.NewRespDataJson(200, bean.ApiData{Code:e.Code, Msg:e.Message})\n}\n"

    // UserModel
    files["app/model/UserModel.go"] = "package model\n\ntype UserModel struct {\n    Table    string  `db:\"user\" json:\"-\"`\n    Connect  string `db:\"default\" json:\"-\"`\n\n    Id       int64 `name:\"id\" pk:\"auto\" json:\"id\"`\n    Password string `name:\"password\" json:\"password,omitempty\"`\n    UserName string `name:\"username\" json:\"username\"`\n\n    CreateAt string `name:\"create_at\" auto:\"time,insert\" json:\"create_at\"`\n    UpdateAt string `name:\"update_at\" auto:\"time,update|insert\" json:\"update_at\"`\n\n    Token    string `json:\"token,omitempty\"`\n}\n    "

    // config.json
    files["config/config.json"] = "{\n  \"evn\": \"debug\",\n  \"app\": {\n    \"debug\": true\n  },\n  \"cache\": {\n    \"cache_driver\": \"redis\"\n  },\n  \"database\": {\n    \"default\": {\n      \"driver\": \"mysql\",\n      \"host\": \"localhost\",\n      \"port\": 3306,\n      \"name\": \"anyminisdk\",\n      \"user\": \"root\",\n      \"password\": \"root\"\n    }\n  },\n  \"redis\": {\n    \"host\": \"127.0.0.1:6379\"\n  }\n}"

    // 写入文件
    for filename, content := range files {
        filename = name + "/src/" + filename
        var f *os.File
        var err error
        if checkFileIsExist(filename) {
            f, err = os.OpenFile(filename, os.O_RDWR, os.ModePerm)
            if err != nil {
                panic(err)
            }

        } else {
            f, err = os.Create(filename)
            if err != nil {
                panic(err)
            }
        }

        f.Write([]byte(content))
    }

    log.Print("create success")
}