package util

import (
	"time"

	"github.com/McaxDev/Back/config"
)

// 加载东八区时区
func UtilInit() {
	var err error
	Loc, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		config.SysLog("ERROR", "加载时区失败")
		return
	}
}
