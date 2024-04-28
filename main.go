package main

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/handler"
	"github.com/McaxDev/Back/util"
)

func main() {

	// 初始化配置文件，自动迁移数据库表
	co.ConfigInit()

	// 初始化工具包
	util.UtilInit()

	// 初始化后端
	handler.HandlerInit()

	// 启动后端
	go Backend()

	// 执行定时任务
	go Cron()

	// 监听命令输入
	Ishell()
}
