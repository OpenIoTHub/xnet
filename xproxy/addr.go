package xproxy

import (
	"github.com/MrMcDuck/xnet/xaddr/xurl"
	"github.com/pkg/errors"
	"strings"
)

/*
 Proxy address example
 "http://my-proxy.com:9090"
 "https://my-proxy.com:9090"
 "socks4://my-proxy.com:9090"
 "socks4a://my-proxy.com:9090"
 "socks5://my-proxy.com:9090"
 "ss://my-proxy.com:9090"
*/

func ParseProxyAddr(addr string) (t ProxyType, host string, err error) {
	us, err := xurl.Parse(addr, "")
	us.Scheme = strings.ToLower(us.Scheme)

	if us.Scheme == "http" {
		t = PROXY_TYPE_HTTP
	} else if us.Scheme == "https" {
		t = PROXY_TYPE_HTTPS
	} else if us.Scheme == "socks4" {
		t = PROXY_TYPE_SOCKS4
	} else if us.Scheme == "socks4a" {
		t = PROXY_TYPE_SOCKS4A
	} else if us.Scheme == "socks5" {
		t = PROXY_TYPE_SOCKS5
	} else if us.Scheme == "ss" {
		t = PROXY_TYPE_SHADOWSOCKS
	} else {
		return t, "", errors.New("Unknow proxy scheme: '" + us.Scheme + "'")
	}

	return t, us.Host.String(), nil
}
