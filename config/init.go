package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 用于操作数据库的变量
var DB *gorm.DB

// 对服务器配置进行初始化
func ConfigInit() {
	var err error
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
