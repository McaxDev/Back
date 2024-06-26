package util

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 读取application/json的请求体并返回gin.H
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

// 清理过期的键值对的值
func ClearExpired[K comparable, V any](timePos func(V) time.Time, maps ...map[K]V) {
	now := time.Now()
	for _, themap := range maps {
		for key, value := range themap {
			if now.After(timePos(value)) {
				delete(themap, key)
			}
		}
	}
}

// 清理过期键值对，但是值是time.Time类型
func ClearExpDefault[K comparable](themap ...map[K]time.Time) {
	ClearExpired(func(t time.Time) time.Time { return t }, themap...)
}
