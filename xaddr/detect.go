package xaddr

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"net"
	"strings"
	"time"
	"github.com/MrMcDuck/xdsa/xstring"
	"github.com/MrMcDuck/xnet/xhttp/client"
	"github.com/MrMcDuck/xnet/xprobe/xonline"
)

// Get all my local IPs
func GetLanIps() ([]net.IP, error) {
	var result = make([]net.IP, 0)

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		if strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "w-") {
			continue
		}

		addrs, err := iface.Addrs()
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

			if ip == nil || ip.IsLoopback() {
				continue
			}

			t := CheckIPString(ip.String())
			if t == IPv4_LAN || t == IPv6_LAN {
				result = append(result, ip)
			}
		}
	}

	return result, nil
}

func GetWanIpOnline() (net.IP, error) {
	var ipstr string
	var exist bool
	var firstCheckWanOnline = true

	// 优先方法
	endpoints := []string{
		"http://whatismyip.akamai.com",
		"http://ident.me",
		"http://myip.dnsomatic.com",
		"http://icanhazip.com"}
	eps := xstring.Shuffle(endpoints)
	for _, url := range eps {
		resp, _, err := client.Get(url, "", time.Second*3, true)
		if err != nil {
			if firstCheckWanOnline {
				if !xonline.IsWanOnline() {
					return nil, errors.New("Can't get WAN ip because of internet offline ")
				}
				firstCheckWanOnline = false
			}
			continue
		}
		ipstr, _ = client.ReadBodyString(resp)
		ipstr = strings.Trim(ipstr, "\r") // icanhazip.com 的返回结果会带换行符
		ipstr = strings.Trim(ipstr, "\n")
		t := CheckIPString(ipstr)
		if t == IPv4_WAN {
			return ParseIP(ipstr)
		}
	}

	// 备用方法
	doc, err := goquery.NewDocument("http://bot.whatismyipaddress.com")
	if err == nil {
		ipstr = doc.Text()
		t := CheckIPString(ipstr)
		if t == IPv4_WAN {
			return ParseIP(ipstr)
		}
	}
	doc, err = goquery.NewDocument("http://network-tools.com/")
	if err == nil {
		ipstr, exist = doc.Find("#field").First().Attr("value")
		if exist {
			t := CheckIPString(ipstr)
			if t == IPv4_WAN {
				return ParseIP(ipstr)
			}
		}
	}

	if !xonline.IsWanOnline() {
		return nil, errors.New("Can't get WAN ip because of internet offline ")
	} else {
		return nil, errors.New("Can't get WAN ip, unknown error")
	}
}
