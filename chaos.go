/*
   @Time : 2020/4/19 1:38 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : chaos
   @Software: GoLand
*/

package chaos

import (
	"context"
	"github.com/offcn-jl/cscf/fake-http"
	"github.com/tencentyun/scf-go-lib/cloudevents/scf"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
)

// HandlerFunc defines the handler used by gin middleware as return value. // fixme 翻译
type HandlerFunc func(*Context)

// HandlersChain defines a HandlerFunc array. // fixme 翻译
type HandlersChain []HandlerFunc

// Last returns the last handler in the chain. ie. the last handler is the main own. // fixme 翻译
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

// Engine is the framework's instance, it contains the muxer, middleware and configuration settings.
// Create an instance of Engine, by using New() or Default() // fixme 翻译
type Engine struct {
	Handlers HandlersChain // 将中间件列表从 router 中提取到 engine 中
}

// New returns a new blank Engine instance without any middleware attached.
// By default the configuration is:
// - RedirectTrailingSlash:  true
// - RedirectFixedPath:      false
// - HandleMethodNotAllowed: false
// - ForwardedByClientIP:    true
// - UseRawPath:             false
// - UnescapePathValues:     true // fixme 翻译
func New() *Engine {
	debugPrintWARNINGNew()
	engine := &Engine{
		//trees:                  make(methodTrees, 0, 9), fixme 未确定是否保留
		//delims:                 render.Delims{Left: "{{", Right: "}}"}, fixme 未确定是否保留
	}
	//engine.RouterGroup.engine = engine fixme 未确定是否保留
	//engine.pool.New = func() interface{} { fixme 未确定是否保留
	//	return engine.allocateContext()
	//}
	return engine
}

// Default returns an Engine instance with the Logger and Recovery middleware already attached. // fixme 翻译
func Default() *Engine {
	debugPrintWARNINGDefault()
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

func (engine *Engine) allocateContext() *Context {
	return &Context{engine: engine}
}

// Use attaches a global middleware to the router. ie. the middleware attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware. // fixme 翻译
func (engine *Engine) Use(middleware ...HandlerFunc) *Engine {
	engine.Handlers = append(engine.Handlers, middleware...) // 将添加中间件的逻辑从 router 中提取到 engine 中进行处理，并且精简掉了 router 中的 GET、 POST 等方法，直接将最终的处理函数作为 USE 的最后一个参数即可
	return engine
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens. // fixme 翻译
// 此函数与 Gin 不同, 没有会返回 err 的步骤, 所以去掉了返回值中的 err, 以及用于 handle err 的 defer
func (engine *Engine) Run() {
	debugPrint("Start Running.")
	cloudfunction.Start(engine.start)
}

// 基于 scf 实现 handleHTTPRequest 的逻辑
func (engine *Engine) start(ctx context.Context, event scf.APIGatewayProxyRequest) (scf.APIGatewayProxyResponse, error) {
	// 初始化上下文
	c := new(Context)

	// 向上下文添加引擎
	c.engine = engine

	// 向上下文添加请求内容
	c.Request = event

	// 初始化引擎
	c.reset()

	// 执行调用链
	c.Next()

	// 如果未设置响应状态码, 则添加默认响应状态码 StatusOK
	if c.Response.StatusCode == 0 {
		c.Response.StatusCode = http.StatusOK
	}

	// 返回响应体给 SCF
	return c.Response, nil
}
