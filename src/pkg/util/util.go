package util

import (
	"fmt"
	"net"
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

func GetOutBoundIP() (ip string, err error) {
	// 使用udp发起网络连接, 这样不需要关注连接是否可通, 随便填一个即可
	return "127.0.0.1", nil
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	// fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}
