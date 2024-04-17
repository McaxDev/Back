package util

import (
	"time"

	co "github.com/McaxDev/Back/config"
)

// 创建时区
var Loc *time.Location

// 加载东八区时区
func init() {
	var err error
	Loc, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		co.SysLog("ERROR", "加载时区失败")
		return
	}
}

// 执行定期任务
func Schedule(second int, operations ...func()) {
	ticker := time.NewTicker(time.Duration(second) * time.Minute)
	for {
		<-ticker.C
		for _, operation := range operations {
			operation()
		}
	}
}
