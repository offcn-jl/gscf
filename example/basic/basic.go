/*
   @Time : 2020/4/20 7:05 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : basic
   @Software: GoLand
*/

package main

import (
	"fmt"
	"github.com/offcn-jl/chaos-go-scf"           // 框架主体
	"github.com/offcn-jl/chaos-go-scf/fake-http" // 假 http 包, fork http 包后移除了大量使用不到的功能
)

func main() {
	r := chaos.Default() // 将 gin 中初始化 engine 的 gin.Default() 修改为 chaos.Default() 即可
	r.Use(basicHandler)  // 由于移除了 router 逻辑，所以直接将 handler 作为 engine.Use 的最后一个参数传给 engine 即可
	_ = r.Run()
}

func basicHandler(c *chaos.Context) {
	fmt.Println("basicHandler is running!")
	count := 408

	// 设置 header
	c.Header("header-1", "hello")
	c.Header("Header-2", "world")

	// 返回数据
	c.JSON(http.StatusOK, chaos.H{"Count": count})
}
