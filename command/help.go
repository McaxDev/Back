package command

import (
	"github.com/spf13/cobra"
)

var Help = &cobra.Command{
	Use:   "help",
	Short: "获取关于这个程序的帮助。",
	Run:   help,
}

func help(cmd *cobra.Command, args []string) {
}
