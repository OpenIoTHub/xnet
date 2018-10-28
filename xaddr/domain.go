package xaddr

//github.com/joeguo/tldextract
//github.com/weppos/publicsuffix-go
//判断是不是公共域名后缀
//判断两个域名是不是同一个所有人，比如news.baidu和www.baidu.com就是同一个所有者

import (
	"github.com/joeguo/tldextract"
	"github.com/pkg/errors"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"os"
	"time"
	"github.com/smcduck/xsys/xenvvar"
	"github.com/smcduck/xsys/xfs"
	"github.com/liamcurry/domains"
	"github.com/domainr/whois"
)

type Domain struct {
	TLD        string // "com" | "com.cn"
	SLD_ROOT   string // "baidu"
	TRD_SUB    string // "www"
	SiteDomain string // "baidu.com"
}

var extractor *tldextract.TLDExtract = nil

// 当缓存文件不存在，或者修改时间在24小时以前，则重新下载缓存文件
func updateLtdListAndExtractor() error {
	var err error
	home, err := xenvvar.GetHomeDir()
	if err != nil {
		return err
	}
	cacheFilename := home + "/.ltd.suffix.cache.dat"
	pi, err := xfs.GetPathInfo(cacheFilename)
	if err != nil {
		extractor = nil
		return err
	}
	if !pi.Exist {
		extractor = nil
	}
	if pi.Exist && time.Now().Sub(pi.ModifiedTime).Hours() > 24 {
		os.Remove(cacheFilename)
		extractor = nil
	}
	if extractor == nil {
		extractor, err = tldextract.New(cacheFilename, false)
		if err != nil {
			return err
		}
	}
	return nil
}

// NOTICE
// 优点: 从权威网站下载TLD列表，判断结果准确
// 缺点: 初始化或者更新时必须在线工作，下载期间接口响应慢
func ParseONLINE(domain string) (*Domain, error) {
	var result Domain

	if err := updateLtdListAndExtractor(); err != nil {
		return nil, err
	}

	if extractor != nil {
		data := extractor.Extract(domain)
		if data != nil && len(data.Tld) > 0 && len(data.Root) > 0 {
			result.TLD = data.Tld
			result.SLD_ROOT = data.Root
			result.TRD_SUB = data.Sub
			result.SiteDomain = result.SLD_ROOT + "." + result.TLD
			return &result, nil
		} else {
			return nil, errors.New("Illegal domain")
		}
	}
	return nil, errors.New("Nil extractor")
}

// NOTICE
// 优点: 响应快，可离线工作
// 缺点: TLD列表固化在代码中，请定期更新库以使判断结果尽可能准确
func ParseOFFLINE(domain string) (*Domain, error) {
	var result Domain

	var fo publicsuffix.FindOptions
	fo.DefaultRule = nil
	fo.IgnorePrivate = false
	dm, err := publicsuffix.ParseFromListWithOptions(publicsuffix.DefaultList, domain, &fo)
	if err != nil {
		return nil, err
	}
	result.SLD_ROOT = dm.SLD
	result.TLD = dm.TLD
	result.TRD_SUB = dm.TRD
	result.SiteDomain = dm.SLD + "." + dm.TLD
	return &result, nil
}

func IsDomainONLINE(domain string) bool {
	_, err := ParseONLINE(domain)
	if err == nil {
		return true
	} else {
		return false
	}
}

func IsDomainOFFLINE(domain string) bool {
	_, err := ParseOFFLINE(domain)
	if err == nil {
		return true
	} else {
		return false
	}
}

// Cloned from github.com/domainr/whois
// Whois response represents a whois response from a server.
type Whois struct {
	// Query and Host are copied from the Request.
	Query string
	Host  string

	// FetchedAt is the date and time the response was fetched from the server.
	FetchedAt time.Time

	// MediaType and Charset hold the MIME-type and character set of the response body.
	MediaType string
	Charset   string

	// Body contains the raw bytes of the network response (minus HTTP headers).
	Body []byte
}

// Principle: WHOIS information of domains which are not taken include "No match".
func IsRegistrable(domain string) bool {
	c := domains.NewChecker()
	return !c.IsTaken(domain)
}

func GetWhois(domain string) (*Whois, error) {
	request, err := whois.NewRequest(domain)
	if err != nil {
		return nil, err
	}
	response, err := whois.DefaultClient.Fetch(request)
	if err != nil {
		return nil, err
	}
	w := Whois{
		Query:     response.Query,
		Host:      response.Host,
		FetchedAt: response.FetchedAt,
		MediaType: response.MediaType,
		Charset:   response.Charset,
		Body:      response.Body,
	}
	return &w, nil
}
