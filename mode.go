/*
   @Time : 2020/4/19 1:52 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : mode
   @Software: GoLand
*/

package gin

import (
	"io"
	"os"
)

// EnvChaosMode 指定配置 Chaos 运行模式的环境变量名称
const EnvChaosMode = "CHAOS_MODE"

const (
	// DebugMode 表示 chaos 模式是 debug
	DebugMode = "debug"
	// ReleaseMode 表示 chaos 模式是 release
	ReleaseMode = "release"
	// TestMode 表示 chaos 模式是 test
	TestMode = "test"
)
const (
	debugCode = iota
	releaseCode
	testCode
)

// DefaultWriter 是 Chaos 用于调试输出和中间件输出的 io.Writer , 例如 Logger() 或 Recovery()
// 注意 : Logger 和 Recovery 都提供了配置其 io.Writer 都自定义方法。
// 要在Windows中支持着色，请使用:
// 		import "github.com/mattn/go-colorable"
// 		gin.DefaultWriter = colorable.NewColorableStdout()
var DefaultWriter io.Writer = os.Stdout

// DefaultErrorWriter is the default io.Writer used by Gin to debug errors fixme 翻译
var DefaultErrorWriter io.Writer = os.Stderr

var chaosMode = debugCode
var modeName = DebugMode

func init() {
	mode := os.Getenv(EnvChaosMode)
	SetMode(mode)
}

// SetMode 根据输入字符串设置 chaos 模式
func SetMode(value string) {
	switch value {
	case DebugMode, "":
		chaosMode = debugCode
	case ReleaseMode:
		chaosMode = releaseCode
	case TestMode:
		chaosMode = testCode
	default:
		panic("chaos mode unknown: " + value)
	}
	if value == "" {
		value = DebugMode
	}
	modeName = value
}
