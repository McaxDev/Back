package util

import (
	"context"
	"time"
)

// 创建一个在指定秒数后超时的context
func Timeout(second int) (context.Context, func()) {
	ctx, canc := context.WithTimeout(context.Background(), time.Minute)
	return ctx, canc
}
