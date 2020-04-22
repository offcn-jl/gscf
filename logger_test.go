/*
   @Time : 2020/4/21 8:56 上午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : logger_test
   @Software: GoLand
*/

package chaos

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/offcn-jl/chaos-go-scf/fake-http"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	SetMode(TestMode)
}

func TestLogger(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithWriter(buffer), func(c *Context) {
		// 一般情况下, 404 状态码由 Api 网关直接返回, 只有某些极特殊的情况是会由 handler 返回, 比如找不到所需资源等
		if c.Request.Path == "/notfound" {
			c.Status(http.StatusNotFound)
		}
	})

	performRequest(router, "GET", "/example?a=100")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
	assert.Contains(t, buffer.String(), "a=100")

	buffer.Reset()
	performRequest(router, "POST", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "POST")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "PUT", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PUT")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "DELETE", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "DELETE")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "PATCH", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PATCH")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "HEAD", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "HEAD")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "OPTIONS", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "OPTIONS")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "GET", "/notfound")
	assert.Contains(t, buffer.String(), "404")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/notfound")
}

func TestLoggerWithConfig(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: buffer}), func(c *Context) {
		// 一般情况下, 404 状态码由 Api 网关直接返回, 只有某些极特殊的情况是会由 handler 返回, 比如找不到所需资源等
		if c.Request.Path == "/notfound" {
			c.Status(http.StatusNotFound)
		}
	})

	performRequest(router, "GET", "/example?a=100")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
	assert.Contains(t, buffer.String(), "a=100")

	buffer.Reset()
	performRequest(router, "POST", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "POST")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "PUT", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PUT")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "DELETE", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "DELETE")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "PATCH", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PATCH")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "HEAD", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "HEAD")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "OPTIONS", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "OPTIONS")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "GET", "/notfound")
	assert.Contains(t, buffer.String(), "404")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/notfound")
}

func TestLoggerWithFormatter(t *testing.T) {
	buffer := new(bytes.Buffer)

	d := DefaultWriter
	DefaultWriter = buffer
	defer func() {
		DefaultWriter = d
	}()

	router := New()
	router.Use(LoggerWithFormatter(func(param LogFormatterParams) string {
		return fmt.Sprintf("[FORMATTER TEST] %v | %3d | %13v | %15s | %-7s %s\n%s",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
			param.ErrorMessage,
		)
	}), func(c *Context) {})
	performRequest(router, "GET", "/example?a=100")

	// output test
	assert.Contains(t, buffer.String(), "[FORMATTER TEST]")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
	assert.Contains(t, buffer.String(), "a=100")
}

func TestLoggerWithConfigFormatting(t *testing.T) {
	var gotParam LogFormatterParams
	var gotKeys map[string]interface{}
	buffer := new(bytes.Buffer)

	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{
		Output: buffer,
		Formatter: func(param LogFormatterParams) string {
			// for assert test
			gotParam = param

			return fmt.Sprintf("[FORMATTER TEST] %v | %3d | %13v | %15s | %-7s %s\n%s",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				param.StatusCode,
				param.Latency,
				param.ClientIP,
				param.Method,
				param.Path,
				param.ErrorMessage,
			)
		},
	}), func(c *Context) {
		// set dummy ClientIP
		c.Request.RequestContext.SourceIP = "20.20.20.20" // c.Request.Header.Set("X-Forwarded-For", "20.20.20.20")
		gotKeys = c.Keys
	})
	performRequest(router, "GET", "/example?a=100")

	// output test
	assert.Contains(t, buffer.String(), "[FORMATTER TEST]")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
	assert.Contains(t, buffer.String(), "a=100")

	// LogFormatterParams test
	//assert.NotNil(t, gotParam.Request)
	assert.NotEmpty(t, gotParam.TimeStamp)
	assert.Equal(t, 200, gotParam.StatusCode)
	assert.NotEmpty(t, gotParam.Latency)
	assert.Equal(t, "20.20.20.20", gotParam.ClientIP)
	assert.Equal(t, "GET", gotParam.Method)
	assert.Equal(t, "/example?a=100", gotParam.Path)
	assert.Empty(t, gotParam.ErrorMessage)
	assert.Equal(t, gotKeys, gotParam.Keys)

}

func TestDefaultLogFormatter(t *testing.T) {
	timeStamp := time.Unix(1544173902, 0).UTC()

	termFalseParam := LogFormatterParams{
		TimeStamp:    timeStamp,
		StatusCode:   200,
		Latency:      time.Second * 5,
		ClientIP:     "20.20.20.20",
		Method:       "GET",
		Path:         "/",
		ErrorMessage: "",
	}
	termFalseLongDurationParam := LogFormatterParams{
		TimeStamp:    timeStamp,
		StatusCode:   200,
		Latency:      time.Millisecond * 9876543210,
		ClientIP:     "20.20.20.20",
		Method:       "GET",
		Path:         "/",
		ErrorMessage: "",
	}

	assert.Equal(t, "[CHOAS] 2018/12/07 - 09:11:42 | 200 |            5s |     20.20.20.20 | GET     /\n", defaultLogFormatter(termFalseParam))
	assert.Equal(t, "[CHOAS] 2018/12/07 - 09:11:42 | 200 |    2743h29m3s |     20.20.20.20 | GET     /\n", defaultLogFormatter(termFalseLongDurationParam))
}

func TestErrorLogger(t *testing.T) {
	router := New()
	router.Use(ErrorLogger(), func(c *Context) {
		c.Error(errors.New("this is an error")) // nolint: errcheck
	})
	w := performRequest(router, "GET", "/error")
	assert.Equal(t, http.StatusOK, w.StatusCode)
	assert.Equal(t, "{\"error\":\"this is an error\"}", w.Body) // 由于 json 生成的方式与 gin 不同, 所以 expected 比 gin 少一个 \n

	router = New()
	router.Use(ErrorLogger(), func(c *Context) {
		c.AbortWithError(http.StatusUnauthorized, errors.New("no authorized")) // nolint: errcheck
	})
	w = performRequest(router, "GET", "/abort")
	assert.Equal(t, http.StatusUnauthorized, w.StatusCode)
	assert.Equal(t, "{\"error\":\"no authorized\"}", w.Body) // 由于 json 生成的方式与 gin 不同, 所以 expected 比 gin 少一个 \n

	router = New()
	router.Use(ErrorLogger(), func(c *Context) {
		c.Error(errors.New("this is an error")) // nolint: errcheck
		c.String(http.StatusInternalServerError, "hola!")
	})
	w = performRequest(router, "GET", "/print")
	assert.Equal(t, http.StatusInternalServerError, w.StatusCode)
	//assert.Equal(t, "hola!{\"error\":\"this is an error\"}\n", w.Body) // 由于 body 的输出 ( 直接给 response.body 赋值 ) 与 gin 框架 ( 持续向 writer 写入数据 ) 不同, 所以没有实现这种 body 连续输出的效果, 所以重写这个测试
	assert.Equal(t, "{\"error\":\"this is an error\"}", w.Body) // 由于 body 的输出 ( 直接给 response.body 赋值 ) 与 gin 框架 ( 持续向 writer 写入数据 ) 不同, 所以没有实现这种 body 连续输出的效果, 所以重写这个测试
}

func TestLoggerWithWriterSkippingPaths(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithWriter(buffer, "/skipped"), func(c *Context) {})

	performRequest(router, "GET", "/logged")
	assert.Contains(t, buffer.String(), "200")

	buffer.Reset()
	performRequest(router, "GET", "/skipped")
	assert.Contains(t, buffer.String(), "")
}

func TestLoggerWithConfigSkippingPaths(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{
		Output:    buffer,
		SkipPaths: []string{"/skipped"},
	}), func(c *Context) {})

	performRequest(router, "GET", "/logged")
	assert.Contains(t, buffer.String(), "200")

	buffer.Reset()
	performRequest(router, "GET", "/skipped")
	assert.Contains(t, buffer.String(), "")
}
