package xproxy

import "net/http"

func GetEnvProxy() (addr string, err error) {
	return "", nil
}

// 参考代码
// 从环境变量$http_proxy或$HTTP_PROXY中获取HTTP代理地址
func GetTransportFromEnvironment() (transport *http.Transport) {
	transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
	return transport
}

func SetEnvProxy() {
}

// 系统环境代理的开关，请参考 https://github.com/txthinking/brook/blob/master/sysproxy/system_darwin.go
