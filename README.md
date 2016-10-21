# bygo 

bygo is a simple and private go webapi framework,it only support api service.

## intro 
 - **simple :** bygo is only a http handler.
 - **private :** i like something,so bygo hava something,i not like something,bygo hava not something,eg:restfull.

## use 
**insert** 
```
go get github.com/bysir-zl/bygo
```
**create project**
 
use bygo ,u need set ```$GOPATH/bin``` in to ```$PATH```
```
bygo create hello_bygo
```
it will create a simple project named hello_bygo

## feature
 - **router :** if u used Laravel(php framework),u will find bygo's route is similar to laravel,Uh uh ,because i like laravel. :D 
 - **middleware :** u can use router+middleware to perform a lot of functions.eg:encrypt the output,decrypt the input,auth from token,add head ...
 - **IOC :** on Controller , if u need db,cache,or custom data in middleware,u just wirte code like this ```func (p IndexController) Index(r *http.Request, s *http.Response) http.ResponseData``` ,
 when bygo call this func,the params will injection.

## to be continued ...

