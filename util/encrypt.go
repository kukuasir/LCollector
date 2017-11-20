package util

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Encrypt(text string) string {
	hashstr := md5.New()
	hashstr.Write([]byte(text))
	return hex.EncodeToString(hashstr.Sum(nil))
}
