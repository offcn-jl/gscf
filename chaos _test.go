/*
   @Time : 2020/4/20 10:44 上午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : chaos _test
   @Software: GoLand
*/

package chaos

import (
	"fmt"
	"github.com/offcn-jl/chaos-go-scf/fake-http"
	"reflect"
	"testing"
)

func compareFunc(t *testing.T, a, b interface{}) {
	sf1 := reflect.ValueOf(a)
	sf2 := reflect.ValueOf(b)
	if sf1.Pointer() != sf2.Pointer() {
		t.Error("different functions")
	}
}

func TestStart(t *testing.T) {
	r := Default()
	t.Log(len(r.Handlers))
	r.Use(testHandler)
	t.Log(len(r.Handlers))
	c := new(Context)
	t.Log(len(c.handlers))
	c.engine = r
	c.reset()
	t.Log(len(c.handlers))
	c.Next()
	t.Log(c.Response)
}

func testHandler(c *Context) {
	count := 999
	fmt.Print("\ntestHandler running.\n\n")
	// 返回数据
	c.JSON(http.StatusOK, H{"Surplus": count})
}
