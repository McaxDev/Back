package command

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ScanCmd() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(" >> ")
	for scanner.Scan() {
		command := strings.TrimSpace(scanner.Text())
		switch command {
		case "help":
			help()
		case "reload":
			reload()
		case "stop", "exit":
			return
		default:
			fmt.Println("未知的命令：", command)
		}
		fmt.Print(" >> ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "读取标准输出失败："+err.Error())
	}
}
