package utils

import (
	"crypto/md5"
	"encoding/hex"
	_"strings"
	_"fmt"
	"strings"
)

/**
	MD5加密
 */
func Md5(str string) (result string, err error) {
	h := md5.New()
	_, err = h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	result = hex.EncodeToString(cipherStr)
	return
}

func FormatStr(s *string) {
	index := strings.LastIndex(*s, "/")
	length := len(*s)
	if index == length - 1 {	//表示该处是以/结尾，去除/
		temp := *s
		*s = temp[:index]
	}


}
