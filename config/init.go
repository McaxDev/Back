package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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

	// 将程序的执行目录更改为项目根目录
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("读取程序所在路径失败：", err.Error())
	}
	if err := os.Chdir(filepath.Dir(exePath)); err != nil {
		log.Fatal("更改程序基准目录失败：", err.Error())
	}

	// 读取配置文件
	if err := Read(&Config, "config.yaml"); err != nil {
		log.Fatalf("重新加载配置文件失败：%v", err)
	}
	if err := Read(&SrvInfo, "srvinfo.json"); err != nil {
		log.Fatalf("重新加载信息文件失败：%v", err)
	}

	// 读取并自动迁移数据库
	DB, err = gorm.Open(mysql.Open(Config.Sql), &gorm.Config{})
	if err != nil {
		log.Fatalf("读取数据库失败：%v", err)
	}
	DB.AutoMigrate(TableList...)

	// 将Redis配置文件里面的数据库编号从字符串变为int
	dbOfRedis, err := strconv.Atoi(Config.Redis["DB"])
	if err != nil {
		log.Fatalf("Redis数据库编号不正确：%v", err)
	}

	// 初始化Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     Config.Redis["host"] + ":" + Config.Redis["port"],
		Password: Config.Redis["password"],
		DB:       dbOfRedis,
	})
}
