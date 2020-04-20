/*
   @Time : 2020/4/19 1:41 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : context
   @Software: GoLand
*/

package chaos

import (
	"github.com/offcn-jl/chaos-go-scf/fake-http"
	"github.com/offcn-jl/chaos-go-scf/render"
	"github.com/tencentyun/scf-go-lib/cloudevents/scf"
	"math"
	"time"
)

const abortIndex int8 = math.MaxInt8 / 2

// Context 是框架最重要的组成部分。
// 它允许我们在中间件之间传递变量、管理链式操作、例如验证请求的 JSON 数据并且返回 JSON 响应。
type Context struct {
	Request  scf.APIGatewayProxyRequest
	Response scf.APIGatewayProxyResponse // fixme response 也许需要单独实现，因为可能需要实现 http 库中提供的一些 getter 或 setter 方法

	engine *Engine

	// Keys is a key/value pair exclusively for the context of each request.
	Keys map[string]interface{}

	handlers HandlersChain
	index    int8

	// Errors is a list of errors attached to all the handlers/middlewares who used this context.
	Errors errorMsgs
}

/************************************/
/************* 创建上下文 *************/
/********** CONTEXT CREATION ********/
/************************************/

func (c *Context) reset() {
	c.Response.Headers = make(map[string]string) // 空 map 需要初始化后才可以使用
	c.Keys = nil
	c.index = -1
	c.handlers = c.engine.Handlers // 增加这一行, 并配合 chaos.start 中的 c.engine = engine 才能给 ctx 增加 handler
	c.Errors = c.Errors[0:0]
}

/************************************/
/************* 调用链控制 *************/
/*********** FLOW CONTROL ***********/
/************************************/

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in GitHub. fixme 阅读理解
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized.
// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called. fixme 阅读理解
func (c *Context) Abort() {
	c.index = abortIndex
}

// AbortWithStatus calls `Abort()` and writes the headers with the specified status code.
// For example, a failed attempt to authenticate a request could use: context.AbortWithStatus(401). fixme 阅读理解
func (c *Context) AbortWithStatus(code int) {
	c.Status(code)
	c.Abort()
}

/************************************/
/************** 错误管理 *************/
/********* ERROR MANAGEMENT *********/
/************************************/

// Error attaches an error to the current context. The error is pushed to a list of errors.
// It's a good idea to call Error for each error that occurred during the resolution of a request.
// A middleware can be used to collect all the errors and push them to a database together,
// print a log, or append it in the HTTP response.
// Error will panic if err is nil. fixme 阅读理解
func (c *Context) Error(err error) *Error {
	if err == nil {
		panic("err is nil")
	}

	parsedError, ok := err.(*Error)
	if !ok {
		parsedError = &Error{
			Err:  err,
			Type: ErrorTypePrivate,
		}
	}

	c.Errors = append(c.Errors, parsedError)
	return parsedError
}

/************************************/
/************* 元数据管理 *************/
/******** METADATA MANAGEMENT *******/
/************************************/

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false) fixme 阅读理解
func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.Keys[key]
	return
}

/************************************/
/************* 生成响应体 *************/
/******** RESPONSE RENDERING ********/
/************************************/

// bodyAllowedForStatus 是 http 包中的 http.bodyAllowedForStatus 函数的副本
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}

// Status 设置 HTTP 响应代码
func (c *Context) Status(code int) {
	// 将 c.writermem.WriteHeader(code)  提取到此处实现
	if code > 0 && c.Response.StatusCode != code {
		if c.Response.StatusCode != 0 {
			debugPrint("[WARNING] Headers were already written. Wanted to override status code %d with %d", c.Response.StatusCode, code)
		}
		c.Response.StatusCode = code
	}
}

// Header is a intelligent shortcut for c.Writer.Header().Set(key, value).
// It writes a header in the response.
// If value == "", this method removes the header `c.Writer.Header().Del(key)`fixme 阅读理解 并且此处也许需要像gin一样使用 header 的 getter setter 来实现
func (c *Context) Header(key, value string) {
	if value == "" {
		delete(c.Response.Headers, key)
		return
	}
	c.Response.Headers[key] = value
}

// Render 写入响应头并调用 render.Render 来呈现数据
func (c *Context) Render(code int, r render.Render) {
	c.Status(code)

	if !bodyAllowedForStatus(code) {
		r.WriteContentType(&c.Response)
		return
	}

	if err := r.Render(&c.Response); err != nil {
		panic(err)
	}
}

// JSON 将给定的结构序列化为 JSON 类型后添加到响应体中
// 并且设置响应头中的 Content-Type 为 "application/json"
func (c *Context) JSON(code int, obj interface{}) {
	c.Render(code, render.JSON{Data: obj})
}

// fixme 阅读理解 应该是对 context 所需方法对实现
/************************************/
/***** GOLANG.ORG/X/NET/CONTEXT *****/
/************************************/

// Deadline returns the time when work done on behalf of this context
// should be canceled. Deadline returns ok==false when no deadline is
// set. Successive calls to Deadline return the same results. fixme 阅读理解
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done returns a channel that's closed when work done on behalf of this
// context should be canceled. Done may return nil if this context can
// never be canceled. Successive calls to Done return the same value. fixme 阅读理解
func (c *Context) Done() <-chan struct{} {
	return nil
}

// Err returns a non-nil error value after Done is closed,
// successive calls to Err return the same error.
// If Done is not yet closed, Err returns nil.
// If Done is closed, Err returns a non-nil error explaining why:
// Canceled if the context was canceled
// or DeadlineExceeded if the context's deadline passed. fixme 阅读理解
func (c *Context) Err() error {
	return nil
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result. fixme 阅读理解
func (c *Context) Value(key interface{}) interface{} {
	if key == 0 {
		return c.Request
	}
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}
	return nil
}
