package main

import (
	"log"
	"os"
	"path/filepath"

	conf "github.com/McaxDev/Back/config"
	hdlr "github.com/McaxDev/Back/handler"
	mid "github.com/McaxDev/Back/middleWare"
	"github.com/gin-gonic/gin"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("读取程序所在路径失败：", err)
	}
	if err := os.Chdir(filepath.Dir(exePath)); err != nil {
		log.Fatal("更改程序基准目录失败：", err)
	}
	if err := conf.ReadConf(); err != nil {
		log.Fatal("配置文件读取失败：", err)
	}
	if err := conf.ReadDB(); err != nil {
		log.Fatal("读取数据库失败：", err)
	}
	r := gin.Default()
	r.GET("/captcha", mid.GetCaptcha)
	r.GET("/challenge", mid.GetChallenge)
	r.GET("/status", hdlr.Status)
	r.GET("/prompt", hdlr.Prompt)
	r.POST("/rcon", mid.AuthCaptcha, mid.AuthChallenge(conf.Config.RconPwd), hdlr.Listenrcon)
	r.POST("/gpt", mid.AuthCaptcha, hdlr.Gpt)
	r.Run(":8080")
}
