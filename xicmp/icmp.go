package xicmp

import (
	"github.com/smcduck/xnet/xicmp/pinger"
	"fmt"
	"os"
	"time"
	"github.com/smcduck/xapputil/xerror"
)

// UDP ping的原理？
// http://blog.csdn.net/duyiwuer2009/article/details/52538638

// 同时支持ICMP-Ping和UDP-Ping的库
// github.com/tobgu/pingo 代码比较啰嗦，带了个什么HTTP服务器
// github.com/arvinkulagin/pinger 缺少IPv6支持 这个能计算延时
// github.com/tatsushid/go-fastping 星最多，支持v4和v6

// "github.com/Cubox-/libping" 简陋杂乱

//https://github.com/tatsushid/go-fastping ICMP/UDP/RawSocket star最多的ping包
//据说里面同时支持udp和raw socket两种方式，可以尝试找找能否调整发送数据包大小和得到正确的分片反馈
// https://github.com/tobgu/pingo TCP/UDP
// https://github.com/arvinkulagin/pinger ICMP/UDP ping

// https://github.com/paulstuart/ping ICMP

type Pong pinger.Pong

func ICMPPing(host string) (*Pong, error) {
	p, err := pinger.NewPinger("icmp", "0.0.0.0")
	p.SetTimeout(time.Second * 5)
	response, err := p.Ping(host, 0)
	return (*Pong)(&response), err
}

func UDPPing(host string) (*Pong, error) {
	p, err := pinger.NewPinger("udp", "0.0.0.0")
	response, err := p.Ping(host, 0)
	return (*Pong)(&response), err
}

func ICMPPingExROOT(host string, payloadLen int) (*Pong, error) {
	p, err := pinger.NewICMP4Pinger("0.0.0.0")
	if err != nil {
		return nil, err
	}
	response, err := p.Ping(host, payloadLen)
	if err != nil {
		return nil, err
	}
	return (*Pong)(&response), nil
}

func UDPPingEx(host string, payloadLen int) (*Pong, error) {
	p, err := pinger.NewUDP4Pinger("0.0.0.0")
	if err != nil {
		return nil, err
	}
	response, err := p.Ping(host, payloadLen)
	if err != nil {
		return nil, err
	}
	fmt.Println(response.RTT)
	return (*Pong)(&response), nil
}

// http://www.letmecheck.it/mtu-test.php
// http://packetlife.net/blog/2008/aug/18/path-mtu-discovery/ 重要参考

/*
https://baike.baidu.com/item/mtu/508920?fr=aladdin
http://blog.sina.com.cn/s/blog_62a5ba4f01018p6d.html

网络中一些常见链路层协议MTU的缺省数值如下：
FDDI协议：4352字节
以太网（Ethernet）协议：1500字节
PPPoE（ADSL）协议：1492字节
X.25协议（Dial Up/Modem）：576字节
Point-to-Point：4470字节
*/

/*
MacOSX下查找MTU的方法
ping -D -s 1465 baidu.com
PING baidu.com (111.13.101.208): 1465 data bytes
556 bytes from oraybox.lan (192.168.9.1): frag needed and DF set (MTU 1492)
Vr HL TOS  Len   ID Flg  off TTL Pro  cks      Src      Dst
 4  5  00 d505 7a9e   0 0000  40  01 1b5b 192.168.9.169  111.13.101.208

Request timeout for icmp_seq 0
*/

const IcmpHeaderSize = 28

func FindMTU2(host string) (int, error) {
	if _, err := UDPPing(host); err != nil {
		return 0, err
	}
	for i := 1400; i < 4471; i++ {
		for j := 0; j < 3; j++ {
			_, err := UDPPingEx(host, i)
			if err == nil {
				break
			} else {
				if j == 2 {
					return i - 1 + IcmpHeaderSize, nil
				}
			}
		}
	}
	return 0, xerror.New("Can't find correct MTU to %s.", host)
}

type findMTUResult struct {
	MTU int
	err error
}

func FindMTU(host string) (int, error) {
	begin := 500
	end := 5000

	if _, err := UDPPing(host); err != nil {
		return 0, err
	}

	occurs := 1
	testsize := (end-begin)*3 + occurs
	retchan := make(chan findMTUResult, testsize)
	seg := (end - begin) / occurs
	for i := 0; i < occurs; i++ {
		bg := begin + (i * seg)
		ed := bg + seg
		if i == occurs-1 {
			ed = end
		}
		bg -= 1
		fmt.Println(bg, ed)
		go findMTU(host, bg, ed, retchan)
	}

	retcnt := 0
	for i := 0; i < occurs; i++ {
		select {
		case ret := <-retchan:
			retcnt += 1
			fmt.Println(ret.err)
			if ret.MTU > 0 && ret.err == nil {
				return ret.MTU, nil
			}
		}
	}

	return 0, xerror.New("Can't find correct MTU between local and %s.", host)
}

func findMTU(host string, begin, end int, retchan chan findMTUResult) {
	ret := findMTUResult{}

	if begin > end {
		ret.MTU = 0
		ret.err = xerror.New("Begin is bigger than end.")
		retchan <- ret
		return
	}

	for i := begin; i < end; i++ {
		for j := 0; j < 3; j++ {
			_, err := UDPPingEx(host, i)
			if err == nil {
				fmt.Println("haha")
				break
			} else {
				if j == 2 {
					ret.MTU = i - 1 + IcmpHeaderSize
					ret.err = nil
					fmt.Println("haha")
					os.Exit(0)
					retchan <- ret
					return
				}
			}
		}
	}
	ret.MTU = 0
	ret.err = xerror.New("Can't find correct MTU between local and %s.", host)
	retchan <- ret
	return
}
