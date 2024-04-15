package util

import (
	"context"
	"time"
)

func Timeout(second int) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	return ctx
}
