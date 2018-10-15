package xaddr

// http://www.baidu.com/news
// socks5://username[:password]@host:port

import (
	"github.com/goware/urlx"
	"github.com/pkg/errors"
	"net/url"
	"strconv"
	"strings"
	"v2ray.com/core/common/net"
	"github.com/MrMcDuck/xdsa/xstring"
	"github.com/MrMcDuck/xarchive/xmultimedia"
)

// scheme:[//[user:password@]host[:port]][/]path[?query][#fragment]

const (
	/* public schemes */
	SchemeUnknown = iota
	SchemeHttp
	SchemeHttps
	SchemeFtp
	SchemeFtps
	SchemeMailto
	SchemeFile
	SchemeIdap
	SchemeNews
	SchemeGopher
	SchemeTelnet
	SchemeWais
	SchemeNntp
	SchemeData
	SchemeIrc
	SchemeIrcs
	SchemeWorldwind
	SchemeMms
	SchemeSocks4
	SchemeSocks4a
	SchemeSocks5
	SchemeSocks5s
	SchemeSocksHttp
	SchemeSocksHttps
	SchemeShadowsocks
	/* custom schemes */
	SchemeSvn
	SchemeHg
	SchemeGit
	SchemeThunder
	SchemeTencent
	SchemeEd2k
	SchemeMagnet
	SchemeTwitter
)

type Scheme int

// "http://bing.com/" is domain url, "http://bing.com/search" is not domain url
func IsDomain(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		return false
	}

	return (len(u.Path) == 0 || u.Path == "/") && len(u.RawQuery) == 0
}

// Combine absolute path and relative path to get a new absolute path
// If relUrl is absolute url, returns this relUrl
func Join(baseUrl string, relUrl string) (absUrl string, err error) {
	if len(baseUrl) == 0 || len(relUrl) == 0 {
		return "", errors.New("UrlJoin get invalid parameters")
	}
	base, err := url.Parse(baseUrl)
	if err != nil {
		return "", errors.New("baseUrl parse error: " + baseUrl)
	}
	if !base.IsAbs() {
		return "", errors.New("baseUrl is not absolute url: " + baseUrl)
	}
	rel, err := url.Parse(relUrl)
	if err != nil {
		return "", errors.New("relUrl parse error: " + relUrl)
	}
	return base.ResolveReference(rel).String(), nil
}

type UrlHost struct {
	Domain string
	IP     string
	Port   int
}

type UrlAuth struct {
	User        string
	Password    string
	PasswordSet bool
}

func (ua *UrlAuth) String() string {
	res := ""
	if len(ua.User) > 0 {
		res += ua.User
	}
	if len(ua.Password) > 0 {
		res += ":" + ua.Password
	}
	return res
}

type UrlSlice struct {
	Scheme string
	Auth   UrlAuth
	Host   UrlHost
	Path   string
}

func (uh *UrlHost) String() string {
	var result string

	if len(uh.Domain) > 0 {
		result += uh.Domain
	} else {
		result += uh.IP
	}

	if IsValidPort(uh.Port) {
		result += ":" + strconv.FormatInt(int64(uh.Port), 10)
	}
	return result
}

func (us *UrlSlice) String() string {
	res := ""
	if len(us.Scheme) > 0 {
		res += us.Scheme + "://"
	}
	if len(us.Auth.String()) > 0 {
		res += us.Auth.String() + "@"
	}
	res += us.Host.String()
	if len(us.Path) > 0 {
		res += "/" + us.Path
	}
	return res
}

// FIXME
// Domain parse undone
//
// NOTICE
// url.Parse("192.168.1.1:80") reports error because RFC3986 says "192.168.1.1:80" is a invalid url, the correct way is "//192.168.1.1:80".
// In xurl, "192.168.1.1:80" is a valid url because it is used a lot
// Reference: https://github.com/golang/go/issues/19297
func ParseUrlOnline(urlStr string, defaultScheme string) (*UrlSlice, error) {

	u, err := urlx.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	var s UrlSlice

	// 即使urlStr中没有Scheme，urlx默认认为并且返回http的scheme
	// 所以这里要检查确认一下返回的Scheme是否正确
	if xstring.StartWith(urlStr, strings.ToLower(u.Scheme) + "://") {
		s.Scheme = u.Scheme
	} else {
		s.Scheme = defaultScheme
	}
	if u.User != nil {
		s.Auth.User = u.User.Username()
		s.Auth.Password, s.Auth.PasswordSet = u.User.Password()
	}
	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		return nil, err
	}
	s.Host.Domain = host
	ip, port, err := ParseHostAddrOnline(u.Host)
	if err != nil {
		return nil, err
	}
	s.Host.IP = ip.String()
	s.Host.Port = port

	s.Path = u.Path

	/*referer_url := "http://www.google.com/search?q=gateway+oracle+cards+denise+linn&hl=en&client=safari"
	r := refererparser.Parse(referer_url)

	log.Printf("Known:%v", r.Known)
	log.Printf("Referer:%v", r.Referer)
	log.Printf("Medium:%v", r.Medium)
	log.Printf("Search parameter:%v", r.SearchParameter)
	log.Printf("Search term:%v", r.SearchTerm)
	log.Printf("Host:%v", r.URI)*/
	return &s, nil
}

func IsImageUrlOnline(url string) bool {
	_, err := ParseUrlOnline(url, "")
	if err != nil {
		return false
	}
	url = strings.ToLower(url)
	for _, v := range xmultimedia.SuffixsOfImage {
		if xstring.EndWith(url, v) {
			return true
		}
	}
	return false
}

func IsVideoUrl(url string) bool {
	_, err := ParseUrlOnline(url, "")
	if err != nil {
		return false
	}
	url = strings.ToLower(url)
	for _, v := range xmultimedia.SuffixsOfVideo {
		if xstring.EndWith(url, v) {
			return true
		}
	}
	return false
}

func IsAudioUrl(url string) bool {
	_, err := ParseUrlOnline(url, "")
	if err != nil {
		return false
	}
	url = strings.ToLower(url)
	for _, v := range xmultimedia.SuffixsOfAudio {
		if xstring.EndWith(url, v) {
			return true
		}
	}
	return false
}

func LastPath(urlstr string) string {
	u, err := url.Parse(urlstr)
	if err != nil {
		return ""
	}

	idx := strings.LastIndex(u.Path, "/")
	if idx <= 0 || idx == (len(u.Path) - 1) {
		return ""
	}
	return u.Path[idx + 1:]
}

// 删除重复的URL，TODO: 这里以后需要改进
func RemoveDuplicateUrl(urls []string) []string {
	return xstring.RemoveDuplicate(urls)
}