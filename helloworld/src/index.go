package main

import
(
    "app/router"
    "net/http"
    "log"
    "github.com/bysir-zl/bygo/bygo"
    "app/exception"
)

func main() {

    apiHandle := bygo.NewApiHandler()

    apiHandle.ConfigRouter("api", router.Init)
    apiHandle.ConfigExceptHandler(exception.Handler)
    apiHandle.Init()

    http.Handle("/api/", apiHandle);
    http.Handle("/", http.FileServer(http.Dir("./dist")))

    log.Println("server start success")

    err := http.ListenAndServe(":81", nil)

    if err != nil {
        log.Println(err)
    }
}
