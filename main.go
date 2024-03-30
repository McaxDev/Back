package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	cmd "github.com/McaxDev/Back/command"
	co "github.com/McaxDev/Back/config"
	hdlr "github.com/McaxDev/Back/handler"
	"github.com/McaxDev/Back/routine"
	"github.com/spf13/cobra"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("读取程序所在路径失败：", err.Error())
	}
	if err := os.Chdir(filepath.Dir(exePath)); err != nil {
		log.Fatal("更改程序基准目录失败：", err.Error())
	}

	co.Init()
	co.AutoMigrate()

	go routine.Backend()
	go routine.Schedule(10, hdlr.ClearExpiredChallenge)

	rootCmd := &cobra.Command{Use: "axoback"}
	rootCmd.AddCommand(cmd.Reload)
	go func() {
		if err := rootCmd.Execute(); err != nil {
			log.Fatalln(err)
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("程序正在退出...")
}
