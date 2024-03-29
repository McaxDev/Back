package handler

import (
	"time"

	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

func LogToSQL(c *gin.Context) {
	startTime := time.Now()
	c.Next()
	duration := time.Since(startTime)
	level := ""
	status := c.Writer.Status()
	switch {
	case status >= 200 && status < 400:
		level = "INFO"
	case status >= 400 && status < 500:
		level = "WARN"
	case status >= 500 && status < 600:
		level = "ERROR"
	}
	util.LogToSQL(c, level, duration)
}
