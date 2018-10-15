package xproxy

import (
	// "github.com/GameXG/ProxyClient" // http client using socks5 proxy supported not well
	"github.com/MrMcDuck/xdsa/xspeed"
	"github.com/MrMcDuck/xnet/xaddr"
	"github.com/MrMcDuck/xnet/xhttp"
	"github.com/pkg/errors"
	"time"
)

const (
	proxyDetectURL = "http://www.baidu.com"
)

type ProxyQuality struct {
	Type      ProxyType
	Available bool
	Speed     *xspeed.Speed
	Latency   time.Duration
}

func CheckProxy(hostAddr string, t ProxyType) (*ProxyQuality, error) {
	if t == PROXY_TYPE_UNKNOWN {
		return nil, errors.New("Unknown proxy type")
	}
	_, _, err := xaddr.ParseHostAddrOnline(hostAddr)
	if err != nil {
		return nil, err
	}

	var pq ProxyQuality
	pq.Available = false
	var counter *xspeed.SpeedCounter = xspeed.NewCounter(time.Minute)

	if t == PROXY_TYPE_HTTP || t == PROXY_TYPE_HTTPS || t == PROXY_TYPE_SOCKS5 {
		counter.BeginCount()
		resp, _, err := xhttp.Get(proxyDetectURL, t.String()+"://"+hostAddr, time.Second*5, true)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			return &pq, nil
		}
		s, err := xhttp.ReadBodyString(resp)
		if err != nil {
			return nil, err
		}
		if len(s) == 0 {
			return nil, errors.New("Empty content")
		}
		pq.Available = true
		counter.Add(uint64(len(resp.Header) + len(s)))
		pq.Speed, err = counter.Get()
		if err != nil {
			return nil, err
		}
		return &pq, nil
	} else {
		return nil, errors.New(t.String() + " type unsupported for now")
	}
}

/*
func TryProxy(hostAddr string) (available bool, t ProxyType, err error) {
	_, _, err = xhostaddr.ParseAddrString(hostAddr)
	if err != nil {
		return false, PROXY_TYPE_UNKNOWN, err
	}

	available, err = CheckProxy(hostAddr, PROXY_TYPE_HTTP)
	if err == nil && available {
		return true, PROXY_TYPE_HTTP, nil
	}
	available, err = CheckProxy(hostAddr, PROXY_TYPE_HTTPS)
	if err == nil && available {
		return true, PROXY_TYPE_HTTPS, nil
	}
	available, err = CheckProxy(hostAddr, PROXY_TYPE_SOCKS5)
	if err == nil && available {
		return true, PROXY_TYPE_SOCKS5, nil
	}
	return false, PROXY_TYPE_UNKNOWN, nil
}*/
