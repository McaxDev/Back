package main

import (
	"log"
	"os"
	"path/filepath"

	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/handler"
	"github.com/gin-gonic/gin"
)

func main() {

	//将Gin设置为发布版
	gin.SetMode(gin.ReleaseMode)

	//将文件执行路径改为当前路径
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("读取程序所在路径失败：", err.Error())
	}
	if err := os.Chdir(filepath.Dir(exePath)); err != nil {
		log.Fatal("更改程序基准目录失败：", err.Error())
	}

	//初始化配置文件，自动迁移数据库表
	co.ConfigInit()

	//启动后端
	handler.HandlerInit()
	go Backend()

	//执行定时任务
	go Cron()

	//监听命令输入
	Ishell()
}
