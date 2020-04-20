/*
   @Time : 2020/4/20 9:19 上午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : test_helpers
   @Software: GoLand
*/

package chaos

import (
	"github.com/tencentyun/scf-go-lib/cloudevents/scf"
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
func performRequest(r *Engine, headers ...header) scf.APIGatewayProxyResponse {
	c := new(Context)
	c.engine = r

	if c.Request.Headers == nil {
		c.Request.Headers = make(map[string]string) // 空 map 需要初始化后才可以使用 : 排除报错 panic: assignment to entry in nil map [recovered]
	}
	for _, h := range headers {
		c.Request.Headers[h.Key] = h.Value
	}

	c.reset()
	c.Next()
	return c.Response
}
