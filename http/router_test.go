package http

import (
	"testing"
	"lib.com/deepzz0/go-com/log"
)

func TestRouter_Start(t *testing.T) {
	context := NewContext()
	router := NewRouter()
	router.Init("", func(node *RouterNode) {
		node.Fun("/api/index", func(c *Context) {
			c.Resp(NewResponseHtml(200, "123"))
			log.Print("qwe")
		})
		node.Group("v1", func(node *RouterNode) {
			node.Fun("/api/index", func(c *Context) {
				c.Resp(NewResponseHtml(200, "123"))
				log.Print("qwe")
			})
		})
	})

	router.Start("/api/index", &context)
}

func BenchmarkRouter(b *testing.B) {
	b.StopTimer()

	context := NewContext()
	router := NewRouter()
	router.Init("", func(node *RouterNode) {
		node.Fun("/api/index2", func(c *Context) {
			c.Resp(NewResponseHtml(200, "123"))
			//log.Print("qwe")
		})
		node.Fun("/api2/index2", func(c *Context) {
			c.Resp(NewResponseHtml(200, "123"))
			//log.Print("qwe")
		})
		node.Fun("/api/index", func(c *Context) {
			c.Resp(NewResponseHtml(200, "123"))
			//log.Print("qwe")
		})
		node.Group("v1", func(node *RouterNode) {
			node.Fun("/api/index", func(c *Context) {
				c.Resp(NewResponseHtml(200, "123"))
				log.Print("qwe")
			})
		})
	})

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		router.Start("/v1/api/index", &context)
	}
}