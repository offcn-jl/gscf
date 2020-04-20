/*
   @Time : 2020/4/19 3:23 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : json
   @Software: GoLand
*/

package render

import (
	"encoding/json"
	"github.com/tencentyun/scf-go-lib/cloudevents/scf"
)

// JSON contains the given interface object. fixme 翻译
type JSON struct {
	Data interface{}
}

var jsonContentType = "application/json; charset=utf-8"

// Render (JSON) writes data with custom ContentType. fixme 翻译
func (r JSON) Render(w *scf.APIGatewayProxyResponse) (err error) {
	if err = WriteJSON(w, r.Data); err != nil {
		panic(err)
	}
	return
}

// WriteContentType (JSON) writes JSON ContentType. fixme 翻译
func (r JSON) WriteContentType(w *scf.APIGatewayProxyResponse) {
	writeContentType(w, jsonContentType)
}

// WriteJSON 生成 json 字符串后写到 APIGatewayProxyResponse 的 Body 中
// WriteJSON marshals the given interface object and writes it with custom ContentType.
func WriteJSON(w *scf.APIGatewayProxyResponse, obj interface{}) error {
	writeContentType(w, jsonContentType)
	jsonBytes, err := json.Marshal(obj)
	w.Body = string(jsonBytes)
	return err
}
