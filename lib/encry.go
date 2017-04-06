package lib

import (
	"crypto/md5"
	"encoding/hex"
	"hash/crc32"
)

//md5加密
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//crc32算法
func Crc32(s string) uint32 {
	h := crc32.NewIEEE()
	h.Write([]byte(s))
	v := h.Sum32()
	return v
}
