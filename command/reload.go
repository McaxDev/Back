package command

import (
	co "github.com/McaxDev/Back/config"
	"github.com/abiosoft/ishell"
)

var loadConfCmd = ishell.Cmd{
	Name: "reload",
	Help: "重新加载配置文件",
	Func: func(c *ishell.Context) {
		co.LoadConfig()
		c.Println("配置文件已重载完成")
	},
}
