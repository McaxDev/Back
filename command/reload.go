package command

import (
	"fmt"

	co "github.com/McaxDev/Back/config"
)

func reload() {
	co.LoadConfig()
	fmt.Println("配置文件已重新加载")
}
