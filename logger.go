/*
   @Time : 2020/4/19 3:55 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : logger
   @Software: GoLand
*/

package chaos

import (
	"fmt"
	"io"
	"time"
)

// todo 完善 logger 包及其单元测试

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

	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// BodySize is the size of the Response Body
	BodySize int
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
	return fmt.Sprintf("[GIN] %v | %3d | %13v | %15s | %-7s %s\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Path,
		param.ErrorMessage,
	)
}

// Logger 初始化一个写入日志到 chaos.DefaultWriter 的日志中间件
// 默认情况下，chaos.DefaultWriter=os.Stdout
// Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
// By default gin.DefaultWriter = os.Stdout.
func Logger() HandlerFunc {
	return LoggerWithConfig(LoggerConfig{})
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
		/*// Start timer fixme
		start := time.Now()
		//path := c.Request.URL.Path
		//raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := LogFormatterParams{
				Request: c.Request,
				Keys:    c.Keys,
			}

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			param.Path = path

			fmt.Fprint(out, formatter(param))
		}*/
	}
}
