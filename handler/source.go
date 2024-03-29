package handler

import (
	"github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/entity"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

func Source(c *gin.Context) {
	title := c.Query("title")
	var text entity.Text
	err := config.DB.Where("title = ?", title).First(&text).Error
	if err != nil {
		util.Warn(c, 400, "查无此文", err)
		return
	}
	textMap := map[string]interface{}{"text": text.Content}
	util.Info(c, 200, "查有此文", textMap)
}
