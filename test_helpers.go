/*
   @Time : 2020/4/20 9:19 上午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : test_helpers
   @Software: GoLand
*/

package chaos

import (
	"github.com/offcn-jl/cscf/fake-http"
	"github.com/tencentyun/scf-go-lib/cloudevents/scf"
	"strings"
)

// CreateTestContext returns a fresh engine and context for testing purposes
func CreateTestContext() (c *Context, r *Engine) {
	r = New()
	c = r.allocateContext()
	c.reset()
	return
}

type header struct {
	Key   string
	Value string
}

// 从 router_test 中迁移至此, 并进行修改
func performRequest(r *Engine, method, path string, headers ...header) scf.APIGatewayProxyResponse {
	c := new(Context)
	c.engine = r

	c.Request.HTTPMethod = method
	c.Request.Path = path
	if len(strings.Split(path, "?")) > 1 { // 将 path 拆分为 SCF API 网关 Event 的数据结构
		c.Request.Path = strings.Split(path, "?")[0]
		if c.Request.QueryString == nil {
			c.Request.QueryString = make(map[string]string)
		}
		for _, value := range strings.Split(strings.Split(path, "?")[1], "&") {
			if len(strings.Split(value, "=")) > 1 {
				c.Request.QueryString[strings.Split(value, "=")[0]] = strings.Split(value, "=")[1]
			} else {
				c.Request.QueryString[strings.Split(value, "=")[0]] = ""
			}
		}
	}
	c.Status(http.StatusOK)

	if c.Request.Headers == nil {
		c.Request.Headers = make(map[string]string)
	}
	for _, h := range headers {
		c.Request.Headers[h.Key] = h.Value
	}

	c.reset()
	c.Next()
	return c.Response
}
