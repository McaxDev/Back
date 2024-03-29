package handler

import (
	"os"

	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

func Variable(c *gin.Context) {
	variables, err := os.ReadFile("variable.json")
	if err != nil {
		util.Error(c, 500, "文件读取失败", err)
		return
	}
	c.Data(200, "application/json", variables)
}
