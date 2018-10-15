package xprobe

import (
	"net"
	"time"
	"github.com/MrMcDuck/xapputil/xerror"
	"github.com/MrMcDuck/xdsa/xstring"
	"github.com/MrMcDuck/xnet/xaddr/xhostaddr"
	"github.com/MrMcDuck/xnet/xaddr/xport"
)

func Tcping(host string, port int, timeout time.Duration) (opened bool, err error) {
	if !xport.IsValidPort(port) {
		return false, xerror.New("Invalid port " + xstring.ToString(port))
	}

	ip, _, err := xhostaddr.ParseAddrString(host)
	if err != nil {
		return false, err
	}

	conn, err := net.DialTimeout("tcp", ip.String()+":"+xstring.ToString(port), timeout)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err != nil {
		return false, nil // Maybe Closed
	}
	return true, nil // Opened
}
