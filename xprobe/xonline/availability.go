package xonline

import (
	"time"
	"github.com/MrMcDuck/xnet/xhttp/client"
	"github.com/MrMcDuck/xdsa/xstring"
)

// don't use LookupHost (DNS) method to detect WAN status,
// when your computer is connected with router OK but router not working well (such like PPPoE account invalid or expired),
// and nameserver is your router
// you can still get DNS record cache by LookupHost API, but in fact the WAN is offline!
func IsWanOnline() (online bool) {
	domains := []string{
		"http://qq.com",
		"http://baidu.com",
		"http://yahoo.com",
		"http://163.com"}

	ds := xstring.Shuffle(domains)
	for _, domain := range ds {
		_, _, err := client.Get(domain, "", time.Second*3, true)
		if err == nil {
			return true
		}
	}
	return false
}
