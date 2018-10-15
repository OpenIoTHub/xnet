package xmkcp

// TODO
// Close mkcp's internal logging output.

import (
	"context"
	"net"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	mkcp "v2ray.com/core/transport/internet/kcp"
	"fmt"
	"github.com/MrMcDuck/xnet/xaddr"
)

// mkcp server
type Server struct {
	ln      *mkcp.Listener
	accepts chan net.Conn
}

func Listen(listenAddr string) (*Server, error) {
	s := Server{}

	s.accepts = make(chan net.Conn, 1024)

	lnNetIP, lnPort, err := xaddr.ParseHostAddrOnline(listenAddr)
	if err != nil {
		return nil, err
	}
	lnIP := lnNetIP.String()

	config := mkcp.Config{
		Mtu:&mkcp.MTU{Value:1500},
		Tti:&mkcp.TTI{Value:10},
		Congestion:true,
		UplinkCapacity:&mkcp.UplinkCapacity{Value:5},
		DownlinkCapacity:&mkcp.DownlinkCapacity{Value:200},
		ReadBuffer:&mkcp.ReadBuffer{Size:10 * 1024 * 1024},
		WriteBuffer:&mkcp.WriteBuffer{Size:10 * 1024 * 1024},
	}
	fmt.Println(config.GetReceivingBufferSize())
	ctx := internet.ContextWithTransportSettings(context.Background(), &config)
	s.ln, err = mkcp.NewListener(ctx, v2net.ParseAddress(lnIP), v2net.Port(lnPort), func(conn internet.Connection) {
		s.accepts <- conn
	})
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Server) Accept() (net.Conn, error) {
	return <-s.accepts, nil
}

func (s *Server) Close() error {
	return s.ln.Close()
}

func (s *Server) Addr() net.Addr {
	return s.ln.Addr()
}

func (s *Server) GetActiveConnCount() int {
	return s.ln.ActiveConnections()
}
