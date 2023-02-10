package network

import (
	"log"
	"net"
	"strconv"
	"time"
)

func Telnet(ip string, port int) bool {
	address := net.JoinHostPort(ip, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		log.Println(err)
		return false
	} else {
		if conn != nil {
			_ = conn.Close()
			return true
		}
	}
	return false
}
