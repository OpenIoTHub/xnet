package socks5

import (
	"fmt"
	"github.com/armon/go-socks5"
	"github.com/pkg/errors"
)

// TODO
// 增加账号密码支持
// 增加级联支持
// 级联使用自定义协议

type Socks5Server struct {
	addr                 string
	username             string
	password             string
	server               *socks5.Server
	network              string
	cascadingSocks5Proxy string
}

// Create a SOCKS5 server
func newServer(network string, addr, username, password, cascadingSocks5Proxy string) (*Socks5Server, error) {
	s := Socks5Server{network: network, addr: addr, username: username, password: password}
	var err error

	if network == "tcp" {
		conf := &socks5.Config{}
		s.server, err = socks5.New(conf)
		if err != nil {
			return nil, err
		}
	} else if network == "udp" {
		conf := &socks5.Config{}
		s.server, err = socks5.New(conf)
		if err != nil {
			return nil, err
		}
	} else if network == "ray" {
		return nil, errors.New(fmt.Sprintf("Ray unsupported for now"))
	} else {
		return nil, errors.New(fmt.Sprintf("Unsupported network %s", network))
	}
	return &s, nil
}

func NewTcpServer(network string, addr, username, password, cascadingSocks5Proxy string) (*Socks5Server, error) {
	return newServer("tcp", addr, username, password, cascadingSocks5Proxy)
}

func NewUdpServer(network string, addr, username, password, cascadingSocks5Proxy string) (*Socks5Server, error) {
	return newServer("udp", addr, username, password, cascadingSocks5Proxy)
}

func NewRayServer(network string, addr, username, password, cascadingSocks5Proxy string) (*Socks5Server, error) {
	return newServer("ray", addr, username, password, cascadingSocks5Proxy)
}

func (s *Socks5Server) ListenAndServe() {
	if s.network == "tcp" {
		s.server.ListenAndServe("tcp", s.addr)
	} else if s.network == "udp" {
		s.server.ListenAndServe("udp", s.addr)
	} else if s.network == "ray" {
	}
}
