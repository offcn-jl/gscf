/*
   @Time : 2020/4/19 3:56 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : recovery
   @Software: GoLand
*/

package chaos

import (
	"bytes"
	"fmt"
	"github.com/offcn-jl/chaos-go-scf/fake-http"
	"io"
	"io/ioutil"
	"log"
	"runtime"
	"time"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// Recovery 返回一个可以 recover 任何 panic 的中间件, 如果有 panic 则返回状态码 500
// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() HandlerFunc {
	return RecoveryWithWriter(DefaultErrorWriter)
}

// RecoveryWithWriter 返回一个 recover 任何 panic 并且可以指定 writer 的中间件, 如果有 panic 则返回状态码 500
// RecoveryWithWriter returns a middleware for a given writer that recovers from any panics and writes a 500 if there was one.
func RecoveryWithWriter(out io.Writer) HandlerFunc {
	var logger *log.Logger
	if out != nil {
		logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	}
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				if logger != nil {
					stack := stack(3)
					headers := c.Request.Headers
					for idx := range headers {
						if idx == "Authorization" {
							headers[idx] = "*"
						}
					}
					if IsDebugging() {
						logger.Printf("[Recovery] %s panic recovered:\n%s\n%s\n%s",
							timeFormat(time.Now()), fmt.Sprint(c.Request), err, stack)
					} else {
						logger.Printf("[Recovery] %s panic recovered:\n%s\n%s",
							timeFormat(time.Now()), err, stack)
					}
				}

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// stack 负责返回格式良好的调用栈, 直接 fork 自 gin 框架中
// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source 返回第 n 行的一个经过空间裁剪的片段, 直接 fork 自 gin 框架中
// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// 返回函数名 如果可的话，函数包含程序计数器 (PC)
// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// 名称没必要包括包的路径名，因为已经包含了文件名，另外它还有一个中心点 "·"
	// 我们获得的是
	// 	runtime/debug.*T·ptrmethod
	// 需要的是
	//	*T.ptrmethod
	// 此外，包路径可能包含点（例如 code.google.com/ ... ），因此需要首先消除路径前缀
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func timeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}
