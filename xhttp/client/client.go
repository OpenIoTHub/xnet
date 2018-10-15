package client

/*
  其实fasthttp也有client的支持，但是不支持代理、验证等等

  参考资料
  使用socks5代理的demo http://mengqi.info/html/2015/201506062329-socks5-proxy-client-in-golang.html
  sosks4(a)代理的支持，可参考https://github.com/h12w/socks & https://github.com/reusee/httpc/blob/master/httpc.go
*/

import (
	. "github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Params:
// proxyAddr 支持http/https/socks5代理
//
// NOTICE
// 如果url不包含http://，将返回错误
// 如果followRedirect==false而且确实发生了跳转，则返回值的redirectUrl将被填写真实的跳转之后的URL；否则redirectUrl返回空
//
// TODO
// 允许从系统默认全局代理发起请求
func Get(url string, proxyAddr string, timeout time.Duration, followRedirect bool) (response *http.Response, redirectUrl *url.URL, err error) {
	request := New()

	var maxRedirectCount int = 5
	var redirectCount int = 0

	// Set timeout
	request.Timeout(timeout)

	// Set proxy
	proxyAddr = strings.ToLower(proxyAddr)
	if len(proxyAddr) > 0 {
		if strings.Index(proxyAddr, "http://") == 0 || strings.Index(proxyAddr, "https://") == 0 {
			request.Proxy(proxyAddr)
		} else if strings.Index(proxyAddr, "socks5://") == 0 {
			proxyAddr = strings.Replace(proxyAddr, "socks5://", "", 1)
			dialer, err := proxy.SOCKS5("tcp", proxyAddr,
				nil,
				&net.Dialer{
					Timeout:   timeout,
					KeepAlive: timeout,
				},
			)
			if err != nil {
				return nil, nil, err
			}
			request.Transport = &http.Transport{
				Proxy:               nil,
				Dial:                dialer.Dial,
				TLSHandshakeTimeout: 10 * time.Second,
			}
		} else {
			return nil, nil, errors.New("Unsupported proxy address " + proxyAddr)
		}
	}

	// Set redirect
	request.RedirectPolicy(
		func(req Request, via []Request) error {
			redirectUrl = req.URL
			if followRedirect {
				redirectCount++
				if redirectCount > maxRedirectCount {
					return errors.New("Too many redirects") // Too many redirects equals HTTP 310 error, but gorequest library doesn't handle this error in it
				}
				return nil // Do redirection
			} else {
				return http.ErrUseLastResponse
			}
		})

	// Run
	resp, _, errs := request.Get(url).End()

	// Return
	if errs != nil { // Handle error
		// Check invalid url
		tmp := strings.ToLower(url)
		if strings.Index(tmp, "http://") != 0 && strings.Index(tmp, "https://") != 0 {
			return nil, nil, errors.New("url must begin with http:// or https://" + ", but input is " + url)
		}
		// Check too many redirects
		if redirectCount > maxRedirectCount {
			resp.StatusCode = 310
		}
		return resp, redirectUrl, errs[0]
	} else { // OK
		return resp, redirectUrl, nil
	}
}

func ReadBodyBytes(response *http.Response) ([]byte, error) {
	return ioutil.ReadAll(response.Body)
}

func ReadBodyString(response *http.Response) (string, error) {
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// 如果没有发生redirect，也能读到一个URL，只不过是原始请求的那个URL
func ReadRedirectUrl(response *http.Response) (*url.URL, error) {
	if response == nil {
		return nil, errors.New("nil input")
	}
	return response.Request.URL, nil
}

func DownloadBigFile(url string, filename string) error {

	// Create the file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
