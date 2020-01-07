package util

import (
	"net"
	"os"
	"strconv"
)

// DefaultHost returns the default host address.
func DefaultHost() string {
	v := os.Getenv("HOST")
	if len(v) > 0 {
		return v
	}

	ip, _ := FindIP(func(ip net.IP) bool {
		return ip.IsLoopback()
	})
	if ip == nil {
		return "127.0.0.1"
	}
	if ip.To4() != nil {
		return ip.String()
	}
	return "[" + ip.To16().String() + "]"
}

// DefaultPort returns the default port number.
func DefaultPort() uint {
	v := os.Getenv("PORT")
	if len(v) == 0 {
		return 5000
	}

	i, err := strconv.Atoi(v)
	if err != nil || i < 0 || i > 65535 {
		return 5000
	}

	return uint(i)
}
