package main

import (
	co "github.com/McaxDev/Back/config"
	"github.com/abiosoft/ishell"
	"github.com/gin-gonic/gin"
)

// 启动监听命令的的函数
func Ishell() {
	shell := ishell.New()
	shell.AddCmd(&ginCmd)
	shell.AddCmd(&loadConfCmd)
	shell.Run()
}

var ginCmd = ishell.Cmd{
	Name: "ginmode",
	Help: "将Gin切换到release、debug或test模式",
	Func: func(c *ishell.Context) {
		if len(c.Args) == 0 {
			c.Println("请提供模式：debug 或 release 或 test")
			return
		}
		mode := c.Args[0]
		switch mode {
		case "debug", "release", "test":
			if mode == gin.Mode() {
				c.Println("当前的Gin模式已经是" + mode)
			}
			gin.SetMode(gin.DebugMode)
			c.Printf("Gin目前运行在%s模式下\n", mode)
		default:
			c.Println("无效的模式，可用模式：debug, release, test")
		}
	},
}

var loadConfCmd = ishell.Cmd{
	Name: "reload",
	Help: "重新加载配置文件",
	Func: func(c *ishell.Context) {
		co.LoadConfig()
		c.Println("配置文件已重载完成")
	},
}
