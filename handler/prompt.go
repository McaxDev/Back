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
)

func Prompt(c *gin.Context) {
	resp, err := http.Get("http://127.0.0.1:1314/status")
	if err != nil {
		c.String(http.StatusInternalServerError, "调用在线人数API失败")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "转换请求体失败")
		return
	}

	var ol_data map[string]interface{}
	if err := json.Unmarshal(body, &ol_data); err != nil {
		c.String(http.StatusInternalServerError, "读取请求体失败")
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
		c.String(http.StatusInternalServerError, "服务器没有图片")
		return
	}
	defer bgfile.Close()

	img, err := png.Decode(bgfile)
	if err != nil {
		c.String(http.StatusInternalServerError, "读取图片失败")
		return
	}

	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	draw.Draw(newImg, bounds, img, bounds.Min, draw.Src)
	fonturi, ok := config.Conf["mcfont"].(string)
	if !ok {
		c.String(http.StatusInternalServerError, "字体路径格式有误")
		return
	}
	fontBytes, err := os.ReadFile(fonturi)
	if err != nil {
		c.String(http.StatusInternalServerError, "未找到字体")
		return
	}
	f, err := opentype.Parse(fontBytes)
	if err != nil {
		c.String(http.StatusInternalServerError, "字体文件有误")
		return
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    18,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "字体设置失败")
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
	versions, ok := config.Conf["ver"].(map[string]interface{})
	if !ok {
		c.String(http.StatusInternalServerError, "类型断言失败")
		return
	}
	mcver, ok := versions["mainbe"].(string)
	if !ok {
		c.String(http.StatusInternalServerError, "类型断言失败")
		return
	}
	d.DrawString("版本：" + mcver)
	d.Dot = fixed.Point26_6{X: fixed.Int26_6(1000), Y: fixed.Int26_6(5000)}
	d.DrawString("点我启动游戏添加服务器")

	c.Writer.Header().Set("Content-type", "image/png")
	if err := png.Encode(c.Writer, newImg); err != nil {
		c.String(http.StatusInternalServerError, "图片发送失败")
	}
}
