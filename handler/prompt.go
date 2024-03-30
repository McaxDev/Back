package handler

import (
	"encoding/json"
	"image/png"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	co "github.com/McaxDev/Back/config"
	"github.com/McaxDev/Back/util"
)

func Prompt(c *gin.Context) {
	resp, err := http.Get("http://127.0.0.1:1314/status")
	if err != nil {
		util.Error(c, 500, "调用在线人数API失败", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.Error(c, 500, "转换请求体失败", err)
		return
	}

	var ol_data map[string]interface{}
	if err := json.Unmarshal(body, &ol_data); err != nil {
		util.Error(c, 500, "读取请求体失败", err)
		return
	}

	var ol_count string
	if player, ok := ol_data["data"].(map[string]interface{}); ok {
		if temp, ok := player["numplayers"].(string); ok {
			ol_count = temp
		}
	}

	bgfile, err := os.Open("onlineBG.png")
	if err != nil {
		util.Error(c, 500, "服务器没有图片", err)
		return
	}
	defer bgfile.Close()

	img, err := png.Decode(bgfile)
	if err != nil {
		util.Error(c, 500, "读取图片失败", err)
		return
	}

	img, err = util.Draw(10, 10, 10, ol_count, "ffffff", img)
	img, err = util.Draw(10, 50, 10, co.SrvInfo.MainVer, "ffffff", img)

	c.Writer.Header().Set("Content-type", "image/png")
	if err := png.Encode(c.Writer, img); err != nil {
		util.Error(c, 500, "图片发送失败", err)
	}
}
