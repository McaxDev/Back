package handler

import (
	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

// 获取用户的文案
func GetText(c *gin.Context) {

	//从查询字符串读取文案类型和标题
	thetype, title := c.Query("type"), c.Query("title")

	//从数据库查找文案
	var tmp co.Text
	if err := co.DB.First(&tmp, "title = ? AND type = ?", title, thetype).Error; err != nil {
		util.Error(c, 400, "查无此文", err)
		return
	}

	//将文案作为响应体返回给客户端
	util.Info(c, 200, "查有此文", gin.H{
		"author":  tmp.User.ID,
		"content": tmp.Content,
	})
}

// 用户设置文案的handler
func SetText(c *gin.Context) {

	//从JWT里获取用户ID
	userID, err := ReadJwt(c)
	if err != nil {
		util.Error(c, 500, "读取用户JWT失败", err)
		return
	}

	//从查询字符串参数获取文案类型，标题和内容
	thetype, title, text := c.PostForm("type"), c.PostForm("title"), c.PostForm("text")

	//检查文案是否已经存在
	var tmp co.Text
	if err := co.DB.First(&tmp, "type = ? AND title = ?", thetype, title).Error; err == nil {
		util.Error(c, 400, "此文已存在", err)
		return
	}
	tmp.Type, tmp.Title, tmp.Content, tmp.AuthorID = thetype, title, text, userID
	if err := co.DB.Create(&tmp).Error; err != nil {
		util.Error(c, 500, "内容创建失败", err)
		return
	}
	util.Info(c, 200, "内容创建成功", nil)
}
