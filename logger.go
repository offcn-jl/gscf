/*
   @Time : 2020/4/19 3:55 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : logger
   @Software: GoLand
*/

package gin

import (
	"fmt"
	"io"
	"time"
)

// 因为云函数环境不支持显示颜色，所以移除了日志颜色相关的处理逻辑

// LoggerConfig 定义日志中间件的配置
// LoggerConfig defines the config for Logger middleware.
type LoggerConfig struct {
	// 可选 默认值为 chaos.defaultLogFormatter
	// Optional. Default value is gin.defaultLogFormatter
	Formatter LogFormatter

	// Output 是写入日志的接口
	// 可选 默认值为 chaos.DefaultWriter
	// Output is a writer where logs are written.
	// Optional. Default value is gin.DefaultWriter.
	Output io.Writer

	// SkipPaths 是一个不写入日志的 url 路径数组
	// 可选
	// SkipPaths is a url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

// LogFormatter 给出格式化程序访问 LoggerWithFormatter 的渠道
// LogFormatter gives the signature of the formatter function passed to LoggerWithFormatter
type LogFormatter func(params LogFormatterParams) string

// LogFormatterParams 是任何日志被格式化时都会被使用的结构
// LogFormatterParams is the structure any formatter will be handed when time to log comes
type LogFormatterParams struct {
	// 由于移除了 http 包及其相关的 Request 逻辑, 所以此处移除了 Request 结构

	// TimeStamp 标记服务器返回响应后经过的时间
	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode 是返回的 HTTP 状态码
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency 是服务器处理某个请求所花费的时间
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP 等于上下文的 ClientIP
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method 是请求的 HTTP 方法
	// Method is the HTTP method given to the request.
	Method string
	// Path 是请求的 Path
	// Path is a path the client requests.
	Path string
	// ErrorMessage 会在处理请求发生错误时被设置
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// BodySize 是响应体的大小
	// BodySize is the size of the Response Body
	BodySize int
	// Keys 是请在请求上下文中设置的键值对
	// Keys are the keys set on the request's context.
	Keys map[string]interface{}
}

// defaultLogFormatter 是记录器中间件使用的默认日志格式函数
// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param LogFormatterParams) string {
	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("[CHOAS] %v | %3d | %13v | %15s | %-7s %s\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Path,
		param.ErrorMessage,
	)
}

// ErrorLogger 返回任何错误类型的 handlerfunc
// ErrorLogger returns a handlerfunc for any error type.
func ErrorLogger() HandlerFunc {
	return ErrorLoggerT(ErrorTypeAny)
}

// ErrorLoggerT 返回给定错误类型的 handlerfunc
// ErrorLoggerT returns a handlerfunc for a given error type.
func ErrorLoggerT(typ ErrorType) HandlerFunc {
	return func(c *Context) {
		c.Next()
		errors := c.Errors.ByType(typ)
		fmt.Println(errors)
		if len(errors) > 0 {
			c.JSON(-1, errors)
		}
	}
}

// Logger 初始化一个写入日志到 chaos.DefaultWriter 的日志中间件
// 默认情况下，chaos.DefaultWriter=os.Stdout
// Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
// By default gin.DefaultWriter = os.Stdout.
func Logger() HandlerFunc {
	return LoggerWithConfig(LoggerConfig{})
}

// LoggerWithFormatter 使用指定的日志格式函数实例一个日志中间件
// LoggerWithFormatter instance a Logger middleware with the specified log format function.
func LoggerWithFormatter(f LogFormatter) HandlerFunc {
	return LoggerWithConfig(LoggerConfig{
		Formatter: f,
	})
}

// LoggerWithWriter 使用给定的 writer buffer 实例一个日志中间件, 并且可以通过 notlogged 跳过不需要记录的内容、
// 例如: os.Stdout, 一个使用 write mode 打开的文件, 一个 socket 链接...
// LoggerWithWriter instance a Logger middleware with the specified writer buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func LoggerWithWriter(out io.Writer, notlogged ...string) HandlerFunc {
	return LoggerWithConfig(LoggerConfig{
		Output:    out,
		SkipPaths: notlogged,
	})
}

// LoggerWithConfig 使用配置初始化一个日志中间件
// LoggerWithConfig instance a Logger middleware with config.
func LoggerWithConfig(conf LoggerConfig) HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = defaultLogFormatter
	}

	out := conf.Output
	if out == nil {
		out = DefaultWriter
	}

	notlogged := conf.SkipPaths

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *Context) {
		// Start timer
		start := time.Now()
		path := c.Request.Path
		raw := ""
		for key, value := range c.Request.QueryString {
			raw += raw + key + "=" + value + "&"
		}
		if len(c.Request.QueryString) > 0 { // 拼接后的 queryString 会多在结尾出一个 & , 进行删除
			raw = raw[0 : len(raw)-1]
		}

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := LogFormatterParams{
				//Request: c.Request,
				Keys: c.Keys,
			}

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.HTTPMethod      // c.Request.Method
			param.StatusCode = c.Response.StatusCode // c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(ErrorTypePrivate).String()

			param.BodySize = len(c.Request.Body) // c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			param.Path = path

			fmt.Fprint(out, formatter(param))
		}
	}
}
