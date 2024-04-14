package command

import (
	"github.com/abiosoft/ishell"
)

// 启动监听命令的的函数
func Ishell() {
	shell := ishell.New()
	shell.AddCmd(&ginCmd)
	shell.AddCmd(&loadConfCmd)
	shell.Run()
}
