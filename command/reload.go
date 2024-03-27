package command

import (
	"fmt"
	"log"

	conf "github.com/McaxDev/Back/config"
	"github.com/spf13/cobra"
)

var Reload = &cobra.Command{
	Use:   "reload",
	Short: "重新加载配置文件并重载程序。",
	Run:   reload,
}

func reload(cmd *cobra.Command, args []string) {
	if err := conf.Read(conf.Config, "config.yaml"); err != nil {
		log.Fatalf("重新加载配置文件失败：%v", err)
	}
	if err := conf.Read(conf.Info, "info.json"); err != nil {
		log.Fatalf("重新加载信息文件失败：%v", err)
	}
	fmt.Println("配置文件已重新加载")
}
