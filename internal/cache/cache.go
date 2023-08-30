package cache

import (
	"log"
	"time"
)

type CaptchaData map[string]interface{}

var CaptchaCache CaptchaData

func clearCaptchaCache() {
	CaptchaCache = nil
}

var RequestCache map[string]string

func clearRequestCache() {
	RequestCache = nil
}

func clearAllCache() {
	clearCaptchaCache()
	clearRequestCache()
	log.Println("5min 清除缓存")
}

// 每5min清除一次缓存
func init() {
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			clearAllCache()
		}
	}()
}
