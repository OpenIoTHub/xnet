package xrudp

import (
	"github.com/u35s/rudp"
	"net"
)

type RUDPConn struct {
	rconn *rudp.RudpConn
}

func Dial(serverAddr string) (*rudp.RudpConn, error) {
	udpaddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return nil, err
	}
	raddr := net.UDPAddr{IP: udpaddr.IP, Port: udpaddr.Port}

	udpconn, err := net.DialUDP("udp", nil, &raddr)
	if err != nil {
		return nil, err
	}
	rconn := rudp.NewConn(udpconn, rudp.New())
	return rconn, nil
}

func Listen(serverAddr string) (*rudp.RudpListener, error) {
	udpaddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return nil, err
	}
	laddr := net.UDPAddr{IP: udpaddr.IP, Port: udpaddr.Port}
	conn, err := net.ListenUDP("udp", &laddr)
	if err != nil {
		return nil, err
	}
	return rudp.NewListener(conn), nil
}
