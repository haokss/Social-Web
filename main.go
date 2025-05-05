package main

import (
	"os"
	"todo_list/config"
	"todo_list/routes"

	log "github.com/sirupsen/logrus"
)

func main() {

	// 创建文件
	file, err := os.OpenFile(".\\log\\app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Error("Can't Create Log File")
	}
	// 初始化配置
	config.Init()

	// 初始化路由
	r := routes.NewRouter()

	// 静态文件处理，允许访问上传的文件
	r.Static("/uploads", "./uploads")
	_ = r.Run(config.HttpPort)
}
