package handler

/*
func lookup(c *gin.Context) {

	//从jwt里读取用户ID
	userID, err := ReadJwt(c)
	if err != nil {
		util.Error(c, 500, "读取JWT失败", err)
		return
	}

	//查询用户的内容
	var tmp co.Text
	err := co.DB.First(&tmp, "author_id = ?" userID).Error
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
*/
