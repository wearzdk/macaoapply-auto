package randomUA

import (
	"math/rand"
	"strconv"
)

// Get 生成随机UA
func Get() string {
	// 随机Chrome版本号
	chromeVersion := strconv.Itoa(rand.Intn(10)+80) + ".0." + strconv.Itoa(rand.Intn(1000)+1000) + "." + strconv.Itoa(rand.Intn(100))
	// 随机Windows版本号
	windowsVersion := strconv.Itoa(rand.Intn(10)+6) + "." + strconv.Itoa(rand.Intn(10)) + "." + strconv.Itoa(rand.Intn(1000)+1000)
	// 拼接UA
	return "Mozilla/5.0 (Windows NT " + windowsVersion + "; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" + chromeVersion + " Safari/537.36"
}
