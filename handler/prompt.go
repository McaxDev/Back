package handler

import (
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"github.com/McaxDev/Back/config"
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

	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	draw.Draw(newImg, bounds, img, bounds.Min, draw.Src)
	fontBytes, err := os.ReadFile(config.Config.McFont)
	if err != nil {
		util.Error(c, 500, "未找到字体", err)
		return
	}
	f, err := opentype.Parse(fontBytes)
	if err != nil {
		util.Error(c, 500, "字体文件有误", err)
		return
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    18,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		util.Error(c, 500, "字体设置失败", err)
		return
	}

	d := &font.Drawer{
		Dst:  newImg,
		Src:  image.NewUniform(color.RGBA{255, 255, 255, 255}),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.Int26_6(1000), Y: fixed.Int26_6(2000)},
	}

	d.DrawString("当前在线 " + ol_count + " 人")
	d.Dot = fixed.Point26_6{X: fixed.Int26_6(1000), Y: fixed.Int26_6(3500)}
	d.DrawString("版本：" + config.Info.MainVer)

	c.Writer.Header().Set("Content-type", "image/png")
	if err := png.Encode(c.Writer, newImg); err != nil {
		util.Error(c, 500, "图片发送失败", err)
	}
}
