/*
   @Time : 2020/4/20 3:57 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : error
   @Software: GoLand
*/

package chaos

// ErrorType 是在 gin 规范中定义的无符号64位错误代码
// ErrorType is an unsigned 64-bit error code as defined in the gin spec.
type ErrorType uint64

const (
	// ErrorTypePrivate 表示私有错误
	// ErrorTypePrivate indicates a private error.
	ErrorTypePrivate ErrorType = 1 << 0
)

// Error 用来规定一个 error 的规范
// Error represents a error's specification.
type Error struct {
	Err  error
	Type ErrorType
	Meta interface{}
}

type errorMsgs []*Error

var _ error = &Error{}

// Error 实现 error 的接口
// Error implements the error interface.
func (msg Error) Error() string {
	return msg.Err.Error()
}
