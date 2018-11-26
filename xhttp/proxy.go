package xhttp

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"time"
	"crypto/tls"
	"github.com/smcduck/xdsa/xstring"
	"github.com/smcduck/xnet/xaddr"
)

// FIXME 不要新建Transport，而是基于原来的修改
// This proxy config method is different from xhttp/client.go Get()
func SetProxy(client *http.Client, proxyUrlString string) error {
	us, err := xaddr.ParseUrl(proxyUrlString)
	if err != nil {
		return err
	}
	if us.Scheme != "http" && us.Scheme != "https" && us.Scheme != "socks5" && us.Scheme != "socks5s" {
		return errors.New(fmt.Sprintf("Unsupported scheme \"%s\"", us.Scheme))
	}

	if xstring.StartWith(us.Scheme, "http") { // Http proxy config
		client.Transport, err = BuildHttpProxyTransport(us.Host.String(), "", "")
		if err != nil {
			return err
		}
	} else { // Socks5 proxy config
		client.Transport, err = BuildSocks5ProxyTransport(us.Host.String(), us.Auth.User, us.Auth.Password)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetInsecureSkipVerify(client *http.Client, skip bool) {
	original := http.DefaultTransport.(*http.Transport)
	if client.Transport != nil {
		original = client.Transport.(*http.Transport)
	}
	if original.TLSClientConfig == nil {
		original.TLSClientConfig = &tls.Config{}
	}
	original.TLSClientConfig.InsecureSkipVerify = skip
	client.Transport = original
}

func SetTimeout(client *http.Client, to time.Duration) {
	client.Timeout = to
}

// Why define myDialer? Cause of golang.org/x/net/proxy need to add DialContext for newest http proxy config apis

type myDialer struct {
	addr   string
	usr    string
	pwd    string
	socks5 proxy.Dialer
}

func (d *myDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	// TODO: golang.org/x/net/proxy need to add DialContext
	return d.Dial(network, addr)
}

func (d *myDialer) Dial(network, addr string) (net.Conn, error) {
	var err error
	if d.socks5 == nil {

		var pauth *proxy.Auth = nil
		auth := proxy.Auth{}
		if len(d.usr) > 0 {
			auth.User = d.usr
			if len(d.pwd) > 0 {
				auth.Password = d.pwd
			}
			pauth = &auth
		}

		d.socks5, err = proxy.SOCKS5("tcp", d.addr, pauth, proxy.Direct)
		if err != nil {
			return nil, err
		}
	}
	return d.socks5.Dial(network, addr)
}

func BuildSocks5ProxyTransport(hostAddr, usr, pwd string) (*http.Transport, error) {
	d := &myDialer{addr: hostAddr, usr: usr, pwd: pwd}
	return &http.Transport{
		DialContext: d.DialContext,
		Dial:        d.Dial,
	}, nil
}

func BuildHttpProxyTransport(hostAddr, usr, pwd string) (*http.Transport, error) {
	proxyURL, err := url.Parse(hostAddr)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		Proxy:               nil,
		//TLSHandshakeTimeout: 10 * time.Second,
	}
	transport.Proxy = http.ProxyURL(proxyURL)
	return transport, nil
}
