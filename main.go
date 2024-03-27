package main

import (
	"log"
	"os"
	"path/filepath"

	cmd "github.com/McaxDev/Back/command"
	conf "github.com/McaxDev/Back/config"
	mid "github.com/McaxDev/Back/middleWare"
	"github.com/McaxDev/Back/routine"
	"github.com/spf13/cobra"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("读取程序所在路径失败：", err)
	}
	if err := os.Chdir(filepath.Dir(exePath)); err != nil {
		log.Fatal("更改程序基准目录失败：", err)
	}
	if err := conf.Read(conf.Config, "config.yaml"); err != nil {
		log.Fatal("配置文件读取失败：", err)
	}
	if err := conf.Read(conf.Info, "info.json"); err != nil {
		log.Fatal("信息读取失败：", err)
	}
	if err := conf.ReadDB(); err != nil {
		log.Fatal("读取数据库失败：", err)
	}

	go routine.Backend()
	go routine.Schedule(10, mid.ClearExpiredChallenge)

	rootCmd := &cobra.Command{Use: "axoback"}
	rootCmd.AddCommand(cmd.Reload)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
