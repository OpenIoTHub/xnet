package xprobe

import (
	"github.com/smcduck/xdsa/xspeed"
	"github.com/smcduck/xnet/xprobe/util/speedtest"
)

// github.com/showwin/speedtest-go 测试下来功能正常，但代码较乱 另外，代码比较清爽但是star很少且没有验证的库 https://github.com/sivel/speedtest/blob/master/speedtest.go

func WanDownloadSpeedTest() (*xspeed.Speed, error) {
	return speedtest.TestDownloadSpeed()
}

func WanUploadSpeedTest() (*xspeed.Speed, error) {
	return speedtest.TestUploadSpeed()
}

// 测试端到端的网速
// https://github.com/blang/speedtest
// https://github.com/DhruvKalaria/SpeedTest
// https://github.com/itimofeev/netspeed
func End2EndUploadSpeedTest() (*xspeed.Speed, error) {
	return nil, nil
}

func End2EndDownloadSpeedTest() (*xspeed.Speed, error) {
	return nil, nil
}
