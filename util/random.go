package util

import (
	"crypto/rand"
	"math/big"
	mtrd "math/rand"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 生成伪随机字符串
func RandStr(n int) string {
	result := make([]byte, n)
	lettersLen := len(letters)
	for i := range result {
		result[i] = letters[mtrd.Intn(lettersLen)]
	}
	return string(result)
}

// 生成真正安全随机的字符串
func RandStrScrure(n int) (string, error) {
	result := make([]byte, n)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		result[i] = letters[num.Int64()]
	}
	return string(result), nil
}
