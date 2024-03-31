package command

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func ginCmd(args []string) {
	subCmd := args[0]
	switch subCmd {
	case "mode":
		mode := args[1]
		switch mode {
		case "debug":
			alterMode(gin.DebugMode)
		case "release":
			alterMode(gin.ReleaseMode)
		default:
			fmt.Println("未知的日志级别：" + mode)
		}
	default:
		fmt.Println("未知的命令：" + subCmd)
	}
}

func alterMode(mode string) {
	if mode == gin.Mode() {
		fmt.Println("当前的GIN模式已经是：" + mode)
	}
	gin.SetMode(mode)
	fmt.Println("已成功将当前的GIN模式设置为：" + mode)
}
