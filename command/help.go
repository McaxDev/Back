package command

import "fmt"

func help() {
	fmt.Println(`
help - 查看命令帮助
reload - 重新加载配置文件
stop - 关闭程序
	`)
}
