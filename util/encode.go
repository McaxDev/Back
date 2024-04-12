package util

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/McaxDev/Back/config"
)

func Encode(origin string, withsalt bool) string {
	if withsalt {
		origin += config.SrvInfo.Salt
	}
	encodedByte := sha256.Sum256([]byte(origin))
	return hex.EncodeToString(encodedByte[:])
}
