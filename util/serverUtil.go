package util

import (
	"fmt"
	"net"
	"time"
)

// 获取服务器Ip
func GetServerIp() (ip string) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
			}
		}
	}

	return
}

/**
获取当前时间
*/
func GetOrderIdTime() (orderId string) {
	currentTime := time.Now().Nanosecond()
	orderId = fmt.Sprintf("%d", currentTime)
	return
}
