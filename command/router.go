package command

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ScanCmd() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := strings.ToLower(strings.TrimSpace(scanner.Text()))
		parts := append(strings.Fields(command), "")
		args := parts[1:]
		switch parts[0] {
		case "help":
			help()
		case "reload":
			reload()
		case "stop", "exit":
			fmt.Println("程序已退出")
			return
		case "log":
			logCmd(args)
		case "gin":
			ginCmd(args)
		default:
			fmt.Println("未知的命令：", command)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "读取标准输出失败："+err.Error())
	}
}
