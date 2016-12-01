package fast

import (
	"github.com/valyala/fasthttp"
	"github.com/bysir-zl/bygo/util"
	"github.com/bysir-zl/fasthttp-routing"
)
// simplest
func Arg2Map(args *fasthttp.Args) map[string]string {
	set := map[string]string{}
	args.VisitAll(func(key []byte, value []byte) {
		set[string(key)] = util.B2S(value)
	})

	return set
}
func Map2Arg(mapper map[string]string) *fasthttp.Args {
	args2 := fasthttp.Args{}
	for k, v := range mapper {
		args2.Add(k, v)
	}
	return &args2
}

func GetPostAll(p *routing.Context) map[string]string {
	return Arg2Map(p.Request.PostArgs())
}