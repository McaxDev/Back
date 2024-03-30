package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	cmd "github.com/McaxDev/Back/command"
	conf "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/entity"
	hdlr "github.com/McaxDev/Back/handler"
	"github.com/McaxDev/Back/routine"
	"github.com/McaxDev/Back/util"
	"github.com/spf13/cobra"
)

func main() {
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

	conf.DB.AutoMigrate(&entity.Log{})
	conf.DB.AutoMigrate(&entity.User{})
	conf.DB.AutoMigrate(&entity.Text{})

	go routine.Backend()
	go routine.Schedule(10, hdlr.ClearExpiredChallenge)

	rootCmd := &cobra.Command{Use: "axoback"}
	rootCmd.AddCommand(cmd.Reload)
	go func() {
		if err := rootCmd.Execute(); err != nil {
			log.Fatalln(err)
		}
	}()

	// 使用 os/signal 包监听系统信号，保持程序运行
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // 阻塞，直到收到终止信号

	// 在这里执行退出前的清理工作
	log.Println("程序正在退出...")
}
