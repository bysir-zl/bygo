package fast

import (
	"reflect"
	"strings"
	"github.com/bysir-zl/fasthttp-routing"
)

func RouterController(r *routing.RouteGroup, path string, controller interface{}, handlers ...routing.Handler) *routing.RouteGroup {
	stru := reflect.ValueOf(controller)
	typ := stru.Type()
	c := r.Group(path)
	// get all func from controller, then add it to router group if sign is 'router.handler'
	for i := stru.NumMethod() - 1; i >= 0; i-- {
		fun := stru.Method(i)
		ifun, ok := fun.Interface().(func(*routing.Context) error)

		if ok {
			name := typ.Method(i).Name
			// when url like this "user/360", but "360" can not ues to function name, u can add "OMIT" prefix on function name, like "user/OMIT360".
			if strings.Index(name, "OMIT") == 0 {
				name = name[4:]
			}
			// skip function
			if strings.Index(name, "SKIP") == 0 {
				continue
			}
			// to lower the initial
			name = strings.ToLower(string(name[0])) + name[1:]
			hs := append([]routing.Handler{func(p *routing.Context) error {
				p.Set("method", name)
				return ifun(p)
			}}, handlers...)
			c.Any("/" + name, hs...)
		}
	}
	return c
}

func Uses(handler routing.Handler, rs ...*routing.RouteGroup) {
	for _, v := range rs {
		v.Use(handler)
	}
}