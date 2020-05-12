/*
   @Time : 2020/4/19 1:52 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : debug
   @Software: GoLand
*/

package gin

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

const chaosSupportMinGoVer = 10

// 如果框架在调试模式下运行，则 IsDebugging 返回 true 。
// 使用 SetMode(gin.ReleaseMode) 禁用调试模式。
func IsDebugging() bool {
	return chaosMode == debugCode
}

func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(DefaultWriter, "[CHAOS-debug] "+format, values...)
	}
}

func getMinVer(v string) (uint64, error) {
	first := strings.IndexByte(v, '.')
	last := strings.LastIndexByte(v, '.')
	if first == last {
		return strconv.ParseUint(v[first+1:], 10, 64)
	}
	return strconv.ParseUint(v[first+1:last], 10, 64)
}

func debugPrintWARNINGDefault() {
	if v, e := getMinVer(runtime.Version()); e == nil && v <= chaosSupportMinGoVer {
		debugPrint(`[WARNING] Now Chaos requires Go 1.10 or later and Go 1.11 will be required soon.

`)
	}
	debugPrint(`[WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

`)
}

func debugPrintWARNINGNew() {
	debugPrint(`[WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export CHAOS_MODE=release
 - using code:	chaos.SetMode(chaos.ReleaseMode)

`)
}

func debugPrintError(err error) {
	if err != nil {
		if IsDebugging() {
			fmt.Fprintf(DefaultErrorWriter, "[CHAOS-debug] [ERROR] %v\n", err)
		}
	}
}
