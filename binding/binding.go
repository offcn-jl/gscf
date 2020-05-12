/*
   @Time : 2020/5/12 1:55 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : binding
   @Software: GoLand
*/

package binding

import (
	"github.com/tencentyun/scf-go-lib/cloudevents/scf"
)

// Binding 描述需要实现的接口，用于绑定请求中的数据，如JSON请求体、查询参数或表单POST
// Binding describes the interface which needs to be implemented for binding the
// data present in the request such as JSON request body, query parameters or
// the form POST.
type Binding interface {
	Name() string
	Bind(scf.APIGatewayProxyRequest, interface{}) error
}

// StructValidator 是需要实现的最小接口
// 此处直接移植了 Gin 的接口定义, 详细描述见下文
// StructValidator is the minimal interface which needs to be implemented in
// order for it to be used as the validator engine for ensuring the correctness
// of the request. Gin provides a default implementation for this using
// https://github.com/go-playground/validator/tree/v8.18.2.
type StructValidator interface {
	// ValidateStruct can receive any kind of type and it should never panic, even if the configuration is not right.
	// If the received type is not a struct, any validation should be skipped and nil must be returned.
	// If the received type is a struct or pointer to a struct, the validation should be performed.
	// If the struct is not valid or the validation itself fails, a descriptive error should be returned.
	// Otherwise nil must be returned.
	ValidateStruct(interface{}) error

	// Engine returns the underlying validator engine which powers the
	// StructValidator implementation.
	Engine() interface{}
}

// Validator 是实现 StructValidator 接口的默认验证器
// Validator is the default validator which implements the StructValidator
// interface. It uses https://github.com/go-playground/validator/tree/v8.18.2
// under the hood.
var Validator StructValidator = &defaultValidator{}

// 这些变量实现了绑定接口, 可用于将请求中的数据绑定到结构体实例
// These implement the Binding interface and can be used to bind the data
// present in the request to struct instances.
var (
	JSON = jsonBinding{}
)

func validate(obj interface{}) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}
