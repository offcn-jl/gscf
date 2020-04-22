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
// Context is the most important part of gin. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
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

// Next 应该确保只在中间件中使用
// 调用 Next 时会开始执行被挂在调用链中的后续 Handlers
// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in GitHub.
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// AbortWithError 会在内部调用 `AbortWithStatus()` 和 `Error()`
// 这个方法可以停止调用链, 并且把指定的 code 写入 HTTP 响应头, 然后把错误推送进 `c.Errors`
// 可以查看 Context.Error() 获取更多细节
// AbortWithError calls `AbortWithStatus()` and `Error()` internally.
// This method stops the chain, writes the status code and pushes the specified error to `c.Errors`.
// See Context.Error() for more details.
func (c *Context) AbortWithError(code int, err error) *Error {
	c.AbortWithStatus(code)
	return c.Error(err)
}

// Abort 可以结束调用链 ( 阻止将要被调用的 handlers 继续被调用 )
// 注意, 这个操作不会停止当前的 handler 继续执行
// 假设您有一个授权中间件来验证当前的请求是否已经授权
// 如果授权失败 ( 例如 : 密码不匹配 ), 则可以调用 Abort 来阻止这个请求上下文中剩余的 handlers
// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized.
// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called.
func (c *Context) Abort() {
	c.index = abortIndex
}

// AbortWithStatus 会调用 `Abort()` 结束调用链, 并将指定的 code 写入 HTTP 响应头
// 例如, 身份验证失败的请求可以使用: context.AbortWithStatus(401)
// AbortWithStatus calls `Abort()` and writes the headers with the specified status code.
// For example, a failed attempt to authenticate a request could use: context.AbortWithStatus(401).
func (c *Context) AbortWithStatus(code int) {
	c.Status(code)
	c.Abort()
}

/************************************/
/************** 错误管理 *************/
/********* ERROR MANAGEMENT *********/
/************************************/

// Error 用来将错误附加到当前上下问中, 附加的错误会被推送到错误列表中
// 对于在解析请求期间发生的错误, 将其传递给 Error 会是一个很好的处理方式
// 中间件可以用来收集所有的错误并将他们一起推送到数据中、打印日志或将其追加到 HTTP 响应中。
// 如果 err 为 nil 时 Error 会返回 panic
// Error attaches an error to the current context. The error is pushed to a list of errors.
// It's a good idea to call Error for each error that occurred during the resolution of a request.
// A middleware can be used to collect all the errors and push them to a database together,
// print a log, or append it in the HTTP response.
// Error will panic if err is nil.
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

// Set 用于在上下文中创建键值对
// 如果 c.Keys 没有被初始化，他还会对其进行初始化
// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
}

// Get 返回给定的 key 对应的 value
// 即: (value, true), 如果给定的值不存在, 则返回 (nil, false)
// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.Keys[key]
	return
}

/************************************/
/************** 请求信息 **************/
/************ INPUT DATA ************/
/************************************/

// Param 可以返回 URL param 的值
// 与 gin 不同, 在本框架中, 它是 c.Request.PathParameters[key] 的语法糖
// Param returns the value of the URL param.
// It is a shortcut for c.Params.ByName(key)
//     router.GET("/user/:id", func(c *gin.Context) {
//         // a GET request to /user/john
//         id := c.Param("id") // id == "john"
//     })
func (c *Context) Param(key string) string {
	return c.Request.PathParameters[key]
}

// ClientIP 在修改后直接返回 SCF Api 网关触发器时间中的 RequestContext.SourceIP
func (c *Context) ClientIP() string {
	return c.Request.RequestContext.SourceIP
}

/************************************/
/************* 生成响应体 *************/
/******** RESPONSE RENDERING ********/
/************************************/

// bodyAllowedForStatus 是 http 包中的 http.bodyAllowedForStatus 函数的副本
// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
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
// Status sets the HTTP response code.
func (c *Context) Status(code int) {
	// 将 c.writermem.WriteHeader(code)  提取到此处实现
	if code > 0 && c.Response.StatusCode != code {
		if c.Response.StatusCode != 0 {
			debugPrint("[WARNING] Headers were already written. Wanted to override status code %d with %d", c.Response.StatusCode, code)
		}
		c.Response.StatusCode = code
	}
}

// Header 可以在响应中写入响应头
// 如果 value == "", 这个方法会删除 key 对应的响应头 fixme 此处也许需要像gin一样使用 header 的 getter setter 来实现
// Header is a intelligent shortcut for c.Writer.Header().Set(key, value).
// It writes a header in the response.
// If value == "", this method removes the header `c.Writer.Header().Del(key)`
func (c *Context) Header(key, value string) {
	if value == "" {
		delete(c.Response.Headers, key)
		return
	}
	c.Response.Headers[key] = value
}

// Render 写入响应头并调用 render.Render 来呈现数据
// Render writes the response headers and calls render.Render to render data.
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

// String 将给定的字符串写入响应体
// String writes the given string into the response body.
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Render(code, render.String{Format: format, Data: values})
}

/************************************/
/*********** 实现上下文接口 ************/
/***** GOLANG.ORG/X/NET/CONTEXT *****/
/************************************/

// Deadline 返回返回此上下文的截止日期
// 当没有设置截止日期时，Deadline 返回 ok == false。 连续调用 Deadline 返回相同的结果
// Deadline returns the time when work done on behalf of this context
// should be canceled. Deadline returns ok==false when no deadline is
// set. Successive calls to Deadline return the same results.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done 返回一个通道，该通道在代表此上下文完成的工作应被取消时关闭。
// 如果无法取消此上下文，则 Done 可能返回 nil 。对 Done 的连续调用返回相同的值。
// Done returns a channel that's closed when work done on behalf of this
// context should be canceled. Done may return nil if this context can
// never be canceled. Successive calls to Done return the same value.
func (c *Context) Done() <-chan struct{} {
	return nil
}

// 在 Done 关闭后返回非 nil 错误值，连续调用 Err 返回相同的错误
// 如果 Done 尚未关闭，Err 将返回 nil
// 如果关闭 Done ，Err 返回一个非 nil 错误，原因是：
// 如果上下文被取消或超过了上下文的截止日期，则会被取消
// Err returns a non-nil error value after Done is closed,
// successive calls to Err return the same error.
// If Done is not yet closed, Err returns nil.
// If Done is closed, Err returns a non-nil error explaining why:
// Canceled if the context was canceled
// or DeadlineExceeded if the context's deadline passed.
func (c *Context) Err() error {
	return nil
}

// Value 返回上下文中的键值对, 如果没有键值对则返回 nil , 上下文中的键值对没有发生改变时连续调用此函数得到的结果相同。
// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
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
