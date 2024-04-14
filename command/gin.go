package command

import (
	"github.com/abiosoft/ishell"
	"github.com/gin-gonic/gin"
)

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
