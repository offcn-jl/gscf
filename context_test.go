/*
   @Time : 2020/4/20 9:10 上午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : context_test.go
   @Software: GoLand
*/

package chaos

import (
	"context"
	"errors"
	"github.com/offcn-jl/cscf/fake-http"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var _ context.Context = &Context{}

func TestContextReset(t *testing.T) { // fixme 实现功能后完善测试
	router := New()
	c := router.allocateContext()
	assert.Equal(t, c.engine, router)

	c.index = 2
	//c.Writer = &responseWriter{ResponseWriter: httptest.NewRecorder()}
	//c.Params = Params{Param{}}
	c.Error(errors.New("test")) // nolint: errcheck
	c.Set("foo", "bar")
	c.reset()

	//assert.False(t, c.IsAborted())
	assert.Nil(t, c.Keys)
	//assert.Nil(t, c.Accepted)
	assert.Len(t, c.Errors, 0)
	assert.Empty(t, c.Errors.Errors())
	assert.Empty(t, c.Errors.ByType(ErrorTypeAny))
	//assert.Len(t, c.Params, 0)
	assert.EqualValues(t, c.index, -1)
	//assert.Equal(t, c.Writer.(*responseWriter), &c.writermem)
}

// 测试上下文中的 handler
func TestContextHandlers(t *testing.T) {
	c, _ := CreateTestContext()
	assert.Nil(t, c.handlers)
	assert.Nil(t, c.handlers.Last())

	c.handlers = HandlersChain{}
	assert.NotNil(t, c.handlers)
	assert.Nil(t, c.handlers.Last())

	f := func(c *Context) {}
	g := func(c *Context) {}

	c.handlers = HandlersChain{f}
	compareFunc(t, f, c.handlers.Last())

	c.handlers = HandlersChain{f, g}
	compareFunc(t, g, c.handlers.Last())
}

// 测试 Status 在设置不同 HTTP 状态码时的表现
func TestStatus(t *testing.T) {
	c, _ := CreateTestContext()

	t.Log(c.Response.StatusCode)

	c.Status(50)
	assert.Equal(t, c.Response.StatusCode, 50)
	t.Log(c.Response.StatusCode)

	c.Status(200)
	assert.Equal(t, c.Response.StatusCode, 200)
	t.Log(c.Response.StatusCode)

	c.Status(500)
	assert.Equal(t, c.Response.StatusCode, 500)
	t.Log(c.Response.StatusCode)

	c.Status(1000)
	assert.Equal(t, c.Response.StatusCode, 1000)
	t.Log(c.Response.StatusCode)

	c.Status(1001)
	assert.Equal(t, c.Response.StatusCode, 1001)
	t.Log(c.Response.StatusCode)
}

// Tests that the response is serialized as JSON fixme 翻译
// and Content-Type is set to application/json
// and special HTML characters are escaped
func TestContextRenderJSON(t *testing.T) {
	c, _ := CreateTestContext()

	c.JSON(http.StatusCreated, H{"foo": "bar", "html": "<b>"})
	c.JSON(http.StatusCreated, H{"foo": "bar", "html": "<b>"})

	assert.Equal(t, http.StatusCreated, c.Response.StatusCode)
	assert.Equal(t, "{\"foo\":\"bar\",\"html\":\"\\u003cb\\u003e\"}", c.Response.Body)
	assert.Equal(t, "application/json; charset=utf-8", c.Response.Headers["Content-Type"])
}

func TestContextHeaders(t *testing.T) {
	c, _ := CreateTestContext()
	c.Header("Content-Type", "text/plain")
	c.Header("X-Custom", "value")

	assert.Equal(t, "text/plain", c.Response.Headers["Content-Type"])
	assert.Equal(t, "value", c.Response.Headers["X-Custom"])

	c.Header("Content-Type", "text/html")
	c.Header("X-Custom", "")

	assert.Equal(t, "text/html", c.Response.Headers["Content-Type"])
	_, exist := c.Response.Headers["X-Custom"]
	assert.False(t, exist)
}

func TestContextGolangContext(t *testing.T) {
	c, _ := CreateTestContext()

	assert.NoError(t, c.Err())
	assert.Nil(t, c.Done())
	ti, ok := c.Deadline()
	assert.Equal(t, ti, time.Time{})
	assert.False(t, ok)
	assert.Equal(t, c.Value(0), c.Request)
	assert.Nil(t, c.Value("foo"))

	c.Set("foo", "bar")
	assert.Equal(t, "bar", c.Value("foo"))
	assert.Nil(t, c.Value(1))
}
