package main

import (
	"log"
	"os"
	"path/filepath"

	cmd "github.com/McaxDev/Back/command"
	conf "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/entity"
	hdlr "github.com/McaxDev/Back/handler"
	"github.com/McaxDev/Back/routine"
	"github.com/McaxDev/Back/util"
	"github.com/spf13/cobra"
)

func main() {

	conf.DB.AutoMigrate(&entity.Log{})
	conf.DB.AutoMigrate(&entity.User{})
	exePath, err := os.Executable()
	if err != nil {
		util.Fatal("读取程序所在路径失败：", err)
	}
	if err := os.Chdir(filepath.Dir(exePath)); err != nil {
		util.Fatal("更改程序基准目录失败：", err)
	}
	if err := conf.Read(&conf.Config, "config.yaml"); err != nil {
		util.Fatal("配置文件读取失败：", err)
	}
	if err := conf.Read(&conf.Info, "info.json"); err != nil {
		util.Fatal("信息读取失败：", err)
	}
	if err := conf.ReadDB(); err != nil {
		util.Fatal("读取数据库失败：", err)
	}

	go routine.Backend()
	go routine.Schedule(10, hdlr.ClearExpiredChallenge)

	rootCmd := &cobra.Command{Use: "axoback"}
	rootCmd.AddCommand(cmd.Reload)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
