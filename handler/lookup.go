package handler

import (
	"os"
	"path/filepath"

	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
)

func lookup(c *gin.Context) {
	username := ReadJwt(c)["username"].(string)
	var tmp co.User
	err := co.DB.Where("username = ?", username).First(&tmp).Error
	if err != nil {
		util.Error(c, 500, "查不到你的信息", err)
		return
	}
	playeruuid := co.PlayerUUID["username"]
	file := playeruuid + ".json"
	server := c.Query("server")
	path := filepath.Join(co.Find(server, "path"), "world/stats/", file)
	playerdata, err := os.ReadFile(path)
	if err != nil {
		util.Error(c, 500, "文件读取失败", err)
		return
	}
	c.Data(200, "application/json", playerdata)
}
