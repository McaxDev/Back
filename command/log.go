package command

import (
	"github.com/abiosoft/ishell"
)

var logCmd = ishell.Cmd{
	Name: "log",
	Help: "Logs a message in the shell",
	Func: func(c *ishell.Context) {
		c.Println("Logging something important!")
	},
}

/*
func logCmd(args []string) {
	subCmd := args[0]
	switch subCmd {
	case "level":
		logLevel := args[1]
		switch logLevel {
		case "debug":
			alterLevel(logrus.DebugLevel)
		case "info":
			alterLevel(logrus.InfoLevel)
		case "warn":
			alterLevel(logrus.WarnLevel)
		case "error":
			alterLevel(logrus.ErrorLevel)
		case "fatal":
			alterLevel(logrus.FatalLevel)
		default:
			fmt.Println("未知的日志级别：" + logLevel)
		}
	default:
		fmt.Println("未知的命令：" + subCmd)
	}
}

func alterLevel(level logrus.Level) {
	levelStr := level.String()
	if level == logrus.GetLevel() {
		fmt.Println("当前的日志级别已经是：" + levelStr)
	}
	logrus.SetLevel(level)
	fmt.Println("已成功将当前的日志级别修改为：" + levelStr)
}
*/
