package util

import (
	"strconv"
	"strings"
)

// 64进制string， 分隔符；，转roomid
func Base64StringToInt64Slice(raw string) (ret []int64) {
	for _, value := range strings.Split(raw, ";") {
		digit, err := strconv.ParseInt(value, 36, 64)
		if err != nil {
			return
		}
		ret = append(ret, digit)
	}
	return
}

func Int64SliceToBase64String(raw []int64) (ret string) {
	var s []string
	for _, value := range raw {
		s = append(s, strconv.FormatInt(value, 36))
	}
	return strings.Join(s, ";")
}
