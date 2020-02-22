package util

import (
	"errors"
	"net"
)

var (
	// ErrNoIP is returned on FindIP cannot find net.IP.
	ErrNoIP = errors.New("No IP addess")
)

// FindIP finds the net.IP in net.InterfaceAddrs().
func FindIP(p func(net.IP) bool) (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if p(ip) {
			return ip, nil
		}
	}

	return nil, ErrNoIP
}
