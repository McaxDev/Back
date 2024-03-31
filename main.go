package main

import (
	"log"
	"os"
	"path/filepath"

	cmd "github.com/McaxDev/Back/command"
	co "github.com/McaxDev/Back/config"
	h "github.com/McaxDev/Back/handler"
	"github.com/McaxDev/Back/routine"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("读取程序所在路径失败：", err.Error())
	}
	if err := os.Chdir(filepath.Dir(exePath)); err != nil {
		log.Fatal("更改程序基准目录失败：", err.Error())
	}

	co.Init()
	co.AutoMigrate()

	go routine.Backend()
	go routine.Schedule(10, h.ClearExpiredChallenge)

	cmd.ScanCmd()
}
