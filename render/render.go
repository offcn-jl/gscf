/*
   @Time : 2020/4/19 3:14 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : render
   @Software: GoLand
*/

package render

import (
	"github.com/tencentyun/scf-go-lib/cloudevents/scf"
)

// Render interface is to be implemented by JSON, XML, HTML, YAML and so on. fixme 翻译
type Render interface {
	// Render writes data with custom ContentType.
	// Render 使用自定义 ContentType 写入数据。 fixme 修正机翻
	Render(*scf.APIGatewayProxyResponse) error
	// WriteContentType writes custom ContentType.
	// WriteContentType 写入自定义 ContentType。 fixme 修正机翻
	WriteContentType(r *scf.APIGatewayProxyResponse)
}

func writeContentType(w *scf.APIGatewayProxyResponse, value string) {
	if val := w.Headers["Content-Type"]; len(val) == 0 {
		w.Headers["Content-Type"] = value
	}
}
