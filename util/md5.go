package util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Encode(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	tempStr := h.Sum(nil)
	return hex.EncodeToString(tempStr)
}
