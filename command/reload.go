package command

import (
	"github.com/spf13/cobra"
)

var Reload = &cobra.Command{
	Use:   "reload",
	Short: "重新加载配置文件并重载程序。",
	Run:   reload,
}

func reload(cmd *cobra.Command, args []string) {
}
