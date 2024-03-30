package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

var queryText = "title = ? AND type = ?"

func GetText(c *gin.Context) {
	thetype, title := c.Query("type"), c.Query("title")
	var tmp co.Text
	err := co.DB.Where(queryText, title, thetype).First(&tmp).Error
	if err != nil {
		util.Warn(c, 400, "查无此文", err)
		return
	}
	textMap := map[string]interface{}{"text": tmp.Content}
	util.Info(c, 200, "查有此文", textMap)
}

func SetText(c *gin.Context) {
	thetype := c.Query("type")
	title, text := c.PostForm("title"), c.PostForm("text")
	var tmp co.Text
	err := co.DB.Where(queryText, title, thetype).First(&tmp).Error
	if err == nil {
		util.Warn(c, 400, "此文已存在", err)
		return
	}
	tmp.Type, tmp.Title, tmp.Content = thetype, title, text
	if err := co.DB.Create(&tmp).Error; err != nil {
		util.Error(c, 500, "内容创建失败", err)
		return
	}
	util.Info(c, 200, "内容创建成功", nil)
}
