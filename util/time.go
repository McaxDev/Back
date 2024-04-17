package util

import (
	"time"
)

// 创建时区
var Loc *time.Location

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
