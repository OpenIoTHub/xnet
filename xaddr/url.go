package xaddr

// http://www.baidu.com/news
// socks5://username[:password]@host:port

import (
	"github.com/goware/urlx"
	"github.com/pkg/errors"
	"github.com/smcduck/xarchive/xmultimedia"
	"github.com/smcduck/xdsa/xstring"
	"net/url"
	"strconv"
	"strings"
	"v2ray.com/core/common/net"
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
func IsAndOnlyDomain(str string) bool {
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

type Path struct {
	Str string
	Dirs []string
	Params map[string][]string
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
	Domain Domain
	Auth   UrlAuth
	Host   UrlHost
	Path   Path
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
	if len(us.Path.Str) > 0 {
		res += "/" + us.Path.Str
	}
	return res
}

// NOTICE
// url.Parse("192.168.1.1:80") reports error because RFC3986 says "192.168.1.1:80" is an invalid url, the correct way is "//192.168.1.1:80".
// In xaddr library, "192.168.1.1:80" is a valid url because it is used a lot
// Reference: https://github.com/golang/go/issues/19297
func ParseUrl(urlStr string) (*UrlSlice, error) {
	u, err := urlx.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	var s UrlSlice

	// if there is NO scheme in input url string, urlx.Parse will give default scheme "http://"
	// so I must check urlx.Parse return scheme
	if xstring.StartWith(strings.ToLower(urlStr), strings.ToLower(u.Scheme) + "://") {
		s.Scheme = u.Scheme
	}
	if u.User != nil {
		s.Auth.User = u.User.Username()
		s.Auth.Password, s.Auth.PasswordSet = u.User.Password()
	}
	if strings.Contains(u.Host, ":") {
		if host, portstr, err := net.SplitHostPort(u.Host); err != nil {
			return nil, err
		} else {
			if portstr != "" {
				if port, err := strconv.Atoi(portstr); err != nil {
					return nil, err
				} else {
					s.Host.Port = port
				}
			}
			s.Host.Domain = host
		}
	} else {
		s.Host.Domain = u.Host
	}
	if s.Host.Domain != "" {
		domain, err := ParseDomain(s.Host.Domain)
		if err != nil {
			return nil, err
		}
		s.Domain = *domain
	}

	s.Path.Str = u.Path
	dirs := strings.Split(u.Path, "/")
	for _, v := range dirs {
		if v == "" {
			continue
		}
		s.Path.Dirs = append(s.Path.Dirs, v)
	}
	s.Path.Params = u.Query()

	return &s, nil
}

func IsImageUrl(url string) bool {
	_, err := ParseUrl(url)
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
	_, err := ParseUrl(url)
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
	_, err := ParseUrl(url)
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