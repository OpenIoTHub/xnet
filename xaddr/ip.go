package xaddr

import (
	"github.com/pkg/errors"
	"net"
)

func ParseIP(s string) (net.IP, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, errors.New("Invalid IP address string '" + s + "'")
	}
	return ip, nil
}
