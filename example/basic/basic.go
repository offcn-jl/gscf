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
	"github.com/offcn-jl/gscf"           // 框架主体
	"github.com/offcn-jl/gscf/fake-http" // 假 http 包, fork http 包后移除了大量 SCF 环境中使用不到的功能
)

func main() {
	r := gin.Default()
	r.Use(basicHandler) // 由于移除了 router 逻辑，所以直接将 handler 作为 engine.Use 的最后一个参数传给 engine 即可
	r.Run()
}

func basicHandler(c *gin.Context) {
	fmt.Println("basicHandler is running!")
	count := 408

	// 设置 header
	c.Header("header-1", "hello")
	c.Header("Header-2", "world")

	// 返回数据
	c.JSON(http.StatusOK, gin.H{"Count": count})
}
