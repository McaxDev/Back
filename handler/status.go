package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/McaxDev/Back/util"
	"github.com/gin-gonic/gin"
	"github.com/mcstatus-io/mcutil/v3"
	"github.com/mcstatus-io/mcutil/v3/response"
)

func Status(c *gin.Context) {
	ctx, canc := context.WithTimeout(context.Background(), time.Second*5)
	defer canc()

	srv := c.Query("srv")

	var resp *response.FullQuery
	var err error

	switch srv {
	case "sc":
		resp, err = mcutil.FullQuery(ctx, "sc.mcax.cn", 25565)
	case "mod":
		resp, err = mcutil.FullQuery(ctx, "mod.mcax.cn", 25565)
	default:
		resp, err = mcutil.FullQuery(ctx, "mcax.cn", 25565)
	}

	if err != nil {
		util.Error(c, 500, "查询失败", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
