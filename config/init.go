package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 用于操作数据库的变量
var DB *gorm.DB

// 对服务器配置进行初始化
func ConfigInit() {
	// 将Gin设置为发布版
	gin.SetMode(gin.ReleaseMode)
	var err error
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("读取程序所在路径失败：", err.Error())
	}
	if err := os.Chdir(filepath.Dir(exePath)); err != nil {
		log.Fatal("更改程序基准目录失败：", err.Error())
	}
	if err := Read(&Config, "config.yaml"); err != nil {
		log.Fatalf("重新加载配置文件失败：%v", err)
	}
	if err := Read(&SrvInfo, "srvinfo.json"); err != nil {
		log.Fatalf("重新加载信息文件失败：%v", err)
	}
	DB, err = gorm.Open(mysql.Open(Config.Sql), &gorm.Config{})
	if err != nil {
		log.Fatalf("读取数据库失败：%v", err)
	}
	DB.AutoMigrate(TableList...)
}
