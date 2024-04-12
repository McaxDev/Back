package main

import (
	"log"
	"os"
	"path/filepath"

	cmd "github.com/McaxDev/Back/command"
	co "github.com/McaxDev/Back/config"
	h "github.com/McaxDev/Back/handler"
	"github.com/McaxDev/Back/routine"
	ut "github.com/McaxDev/Back/util"
)

func main() {

	//将Gin设置为发布版
	//gin.SetMode(gin.ReleaseMode)

	//将文件执行路径改为当前路径
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("读取程序所在路径失败：", err.Error())
	}
	if err := os.Chdir(filepath.Dir(exePath)); err != nil {
		log.Fatal("更改程序基准目录失败：", err.Error())
	}

	//初始化配置文件，自动迁移数据库表
	co.Init()
	co.AutoMigrate()

	//启动后端
	go routine.Backend()

	//执行定时任务
	go routine.Schedule(10,
		ut.ClearExpired(h.Challenges),
		ut.ClearExpired(h.IpTimeMap),
		h.ClearExpiredMailSent,
	)

	//监听命令输入
	cmd.ScanCmd()
}
