package util

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 上传文件到指定路径的函数
func UploadFile(c *gin.Context, key string, filedir ...string) (string, error) {

	// 从请求体获取文件
	file, err := c.FormFile(key)
	if err != nil {
		return "", err
	}

	// 获取文件的路径名和文件名
	filename := uuid.New().String() + filepath.Ext(file.Filename)
	innerDir := filepath.Join(filedir...)

	// 确保上传目录存在
	outerDir := filepath.Join("assets", innerDir)
	if _, err := os.Stat(outerDir); os.IsNotExist(err) {
		if err := os.Mkdir(outerDir, os.ModePerm); err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	// 保存文件
	path := filepath.Join(outerDir, filename)
	if err := c.SaveUploadedFile(file, path); err != nil {
		return "", err
	}

	// 返回内层文件路径
	return filepath.Join(innerDir, filename), nil
}
