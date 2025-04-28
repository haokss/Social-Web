package main

import (
	"todo_list/config"
	"todo_list/routes"
)

func main() {
	// 初始化配置
	config.Init()
	// 初始化事件
	// go service.StartTaskScheduler()
	// 初始化路由
	r := routes.NewRouter()
	// 静态文件处理，允许访问上传的文件
	r.Static("/uploads", "./uploads")
	_ = r.Run(config.HttpPort)
}
