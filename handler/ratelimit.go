package handler

import (
	"fmt"
	"sync"
	"time"

	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 返回限频中间件的工厂函数
func RateLimit(seconds int) gin.HandlerFunc {

	// 创建自由变量
	visits := &struct {
		sync.RWMutex
		lastTime map[string]time.Time
	}{lastTime: make(map[string]time.Time)}

	// 返回闭包
	return func(c *gin.Context) {

		// 获取用户IP
		ip := c.ClientIP()

		// 上锁并获取上次的访问时间
		visits.RLock()
		lastvisit, ok := visits.lastTime[ip]
		visits.RUnlock()

		// 检查用户的访问频率是否超过了限制
		timeLeft := seconds - int(time.Since(lastvisit).Seconds())
		if ok && timeLeft > 0 {
			msg := fmt.Sprintf("你访问的太快了，请 %d 秒后重试。", timeLeft)
			util.Error(c, 400, msg, nil)
			return
		}

		// 上锁并更新用户最新的记录
		visits.Lock()
		visits.lastTime[ip] = time.Now()
		visits.Unlock()
	}
}
