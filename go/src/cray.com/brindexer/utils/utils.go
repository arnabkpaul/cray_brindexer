package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func Md5FromString(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}
