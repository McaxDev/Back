package util

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MapReadResp(res *http.Response) (gin.H, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("响应体读取失败")
	}
	var data gin.H
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, errors.New("JSON反序列化失败")
	}
	return data, nil
}
