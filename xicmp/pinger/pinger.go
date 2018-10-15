package pinger

import (
	"errors"
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"math/rand"
	"net"
	"time"
)

const (
	protocolICMP = 1
)

type Pinger interface {
	Ping(raddr string, payloadLen int) (Pong, error)
	ResetCounter()
	SetTimeout(time.Duration)
}

// Returns new Pinger interface. network argument must be "icmp" or "udp" (ICMP6 and UDP6 support will be added later).
func NewPinger(network, laddr string) (Pinger, error) {
	switch network {
	case "icmp":
		return NewICMP4Pinger(laddr)
	case "udp":
		return NewUDP4Pinger(laddr)
	default:
		return &ICMP4Pinger{}, errors.New("Unknown network " + network)
	}
}

type ICMP4Pinger struct {
	laddr   net.Addr
	id      int
	counter int
	timeout time.Duration
}

// Returns new ICMP4Pinger. laddr is local ip address for listening.
func NewICMP4Pinger(laddr string) (*ICMP4Pinger, error) {
	addr, err := net.ResolveIPAddr("ip4", laddr)
	if err != nil {
		return nil, err
	}
	pinger := ICMP4Pinger{
		laddr:   addr,
		id:      rand.Int() & 0xffff,
		counter: 0,
		timeout: 2 * time.Second,
	}
	return &pinger, nil
}

// Sets ICMP4Pinger counter to 0. Counter increments with each Ping() call.
// Counter value is set to Seq field in Echo-Request.
func (i *ICMP4Pinger) ResetCounter() {
	i.counter = 0
}

// Sets ICMP4Pinger timeout. Timeout is a waiting time for a Echo-Reply from a remote host.
func (i *ICMP4Pinger) SetTimeout(d time.Duration) {
	i.timeout = d
}

// SendEth Echo-Request to remote host and wait Echo-Reply.
// raddr is an address of remote host.
func (i *ICMP4Pinger) Ping(raddr string, payloadLen int) (Pong, error) {
	buf := []byte{}
	if payloadLen > 0 {
		buf = make([]byte, payloadLen)
	}
	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   i.id,
			Seq:  i.counter,
			Data: buf,
		},
	}
	i.counter++
	listener, err := icmp.ListenPacket("ip4:icmp", i.laddr.String())
	if err != nil {
		return Pong{}, err
	}
	defer listener.Close()
	addr, err := net.ResolveIPAddr("ip4", raddr)
	if err != nil {
		return Pong{}, err
	}
	return ping(listener, message, addr, i.timeout)
}

type UDP4Pinger struct {
	laddr   net.Addr
	id      int
	counter int
	timeout time.Duration
}

// Returns new UDP4Pinger. laddr is local ip address for listening.
func NewUDP4Pinger(laddr string) (*UDP4Pinger, error) {
	addr, err := net.ResolveIPAddr("ip4", laddr)
	if err != nil {
		return nil, err
	}
	pinger := UDP4Pinger{
		laddr:   addr,
		id:      rand.Int() & 0xffff,
		counter: 0,
		timeout: 2 * time.Second,
	}
	return &pinger, nil
}

// Sets UDP4Pinger counter to 0. Counter increments with each Ping() call.
// Counter value is set to Seq field in Echo-Request.
func (i *UDP4Pinger) ResetCounter() {
	i.counter = 0
}

// Sets UDP4Pinger timeout. Timeout is a waiting time for a Echo-Reply from a remote host.
func (i *UDP4Pinger) SetTimeout(d time.Duration) {
	i.timeout = d
}

// SendEth Echo-Request to remote host and wait Echo-Reply.
// raddr is an address of remote host.
func (i *UDP4Pinger) Ping(raddr string, payloadLen int) (Pong, error) {
	buf := []byte{}
	if payloadLen > 0 {
		fmt.Println(payloadLen)
		buf = make([]byte, payloadLen)
	} else {
		buf = nil
	}
	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   i.id,
			Seq:  i.counter,
			Data: buf,
		},
	}
	i.counter++
	listener, err := icmp.ListenPacket("udp4", i.laddr.String())
	if err != nil {
		return Pong{}, err
	}
	defer listener.Close()
	addr, err := net.ResolveIPAddr("ip4", raddr)
	if err != nil {
		return Pong{}, err
	}
	return ping(listener, message, &net.UDPAddr{IP: net.ParseIP(addr.String())}, i.timeout)
}

func ping(listener *icmp.PacketConn, message icmp.Message, raddr net.Addr, timeout time.Duration) (Pong, error) {
	data, err := message.Marshal(nil)
	if err != nil {
		return Pong{}, err
	}
	n, err := listener.WriteTo(data, raddr)
	if err != nil {
		return Pong{}, err
	}
	if n != message.Body.Len(0)+4 {
		return Pong{}, errors.New("Write size error")
	}
	now := time.Now()
	done := make(chan Pong)
	errch := make(chan string, 10) // 为什么errch从来都不管用？
	go func() {
		for {
			buf := make([]byte, 10000)
			// bufio
			n, ra, err := listener.ReadFrom(buf)
			if err != nil {
				errch <- err.Error()
				return
			}
			since := time.Since(now)
			input, err := icmp.ParseMessage(protocolICMP, buf[:n])
			if err != nil {
				errch <- err.Error()
				return
			}
			if input.Type != ipv4.ICMPTypeEchoReply {
				continue
			}
			echo := input.Body.(*icmp.Echo)
			pong := Pong{
				RemoteAddr: ra,
				ID:         echo.ID,
				Seq:        echo.Seq,
				Data:       echo.Data,
				Size:       n,
				RTT:        since,
			}
			done <- pong
			return
		}
	}()
	select {
	case pong := <-done:
		return pong, nil
	case errstr := <-errch:
		return Pong{}, errors.New(errstr)
	case <-time.After(timeout):
		return Pong{}, errors.New("Timeout")
	}
}

type Pong struct {
	// IP address of pinged host.
	RemoteAddr net.Addr
	// ICMP ID.
	ID int
	// ICMP sequence number.
	Seq int
	// Content of ICMP data field.
	Data []byte
	// Size of ICMP Echo-Reply.
	Size int
	// Round-trip time.
	RTT time.Duration
}
