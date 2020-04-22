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

	//encoder := json.NewEncoder(w)  // gin 框架原来的生成 json 方式
	//err := encoder.Encode(&obj)  // gin 框架原来的生成 json 方式

	//buffer := new(bytes.Buffer) // 模拟 gin 框架 json 生成方式
	//encoder := json.NewEncoder(buffer)  // 模拟 gin 框架 json 生成方式
	//err := encoder.Encode(&obj) // 模拟 gin 框架 json 生成方式
	//w.Body = buffer.String() // 模拟 gin 框架 json 生成方式, 这样生成出来的 json 会带有一个 \n

	jsonBytes, err := json.Marshal(obj)
	w.Body = string(jsonBytes) // 这样生成的 json 会比 gin 框架少一个 \n , 本质原因是少了 json.encoder 添加的 \n
	return err
}
