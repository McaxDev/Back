package handler

import (
	"os"

	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

func Variable(c *gin.Context) {
	variables, err := os.ReadFile("variable.json")
	if err != nil {
		c.JSON(500, util.Json("文件读取失败", nil))
		return
	}
	c.Data(200, "application/json", variables)
}