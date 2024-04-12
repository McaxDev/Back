package util

import (
	"crypto/sha256"
	"encoding/hex"

	co "github.com/McaxDev/Back/config"
)

func Encode(origin string) string {
	encodedByte := sha256.Sum256([]byte(origin + co.SrvInfo.Salt))
	return hex.EncodeToString(encodedByte[:])
}
