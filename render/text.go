/*
   @Time : 2020/4/21 10:05 上午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : text
   @Software: GoLand
*/

package render

import (
	"fmt"
	"github.com/tencentyun/scf-go-lib/cloudevents/scf"
)

// String 包含给定的接口对象切片及其格式化条件
// String contains the given interface object slice and its format.
type String struct {
	Format string
	Data   []interface{}
}

var plainContentType = "text/plain; charset=utf-8"

// Render (String) 使用自定义 ContentType 写入数据
// Render (String) writes data with custom ContentType.
func (r String) Render(w *scf.APIGatewayProxyResponse) error {
	return WriteString(w, r.Format, r.Data)
}

// WriteContentType (String) 写入 ContentType 到 HTTP 响应头
// WriteContentType (String) writes Plain ContentType.
func (r String) WriteContentType(w *scf.APIGatewayProxyResponse) {
	writeContentType(w, plainContentType)
}

// WriteString 根据格式化条件格式化写入数据到 HTTP 响应体, 并将 ContentType 写入到 HTTP 响应头
// WriteString writes data according to its format and write custom ContentType.
func WriteString(w *scf.APIGatewayProxyResponse, format string, data []interface{}) (err error) {
	writeContentType(w, plainContentType)
	w.Body = fmt.Sprintf(format, data...)
	return
}
