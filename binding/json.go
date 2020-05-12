/*
   @Time : 2020/5/12 1:58 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : json
   @Software: GoLand
*/

package binding

import (
	"encoding/json"
	"fmt"
	"github.com/tencentyun/scf-go-lib/cloudevents/scf"
	"strings"
)

// EnableDecoderUseNumber 用于调用JSON解码器实例上的UseNumber方法
// UseNumber 使解码器将一个数字作为数字而不是 float64 解组到 interface{}
// EnableDecoderUseNumber is used to call the UseNumber method on the JSON
// Decoder instance. UseNumber causes the Decoder to unmarshal a number into an
// interface{} as a Number instead of as a float64.
var EnableDecoderUseNumber = false

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Bind(req scf.APIGatewayProxyRequest, obj interface{}) error {
	if req.Body == "" {
		return fmt.Errorf("invalid request")
	}
	return decodeJSON(req.Body, obj)
}

func decodeJSON(r string, obj interface{}) error {
	decoder := json.NewDecoder(strings.NewReader(r)) // 将此处的 io.Reader 替换为了 strings.NewReader 来读取 SCF 请求体中的 body 字符串
	if EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}
