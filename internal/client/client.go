package client

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"macaoapply-auto/internal/cache"
	"macaoapply-auto/pkg/config"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tidwall/gjson"
)

var client *resty.Client

// 累计请求次数
var requestCount int

// 慢网络检测
var isSlowNetwork = false

func genIss() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		panic(err)
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	iss := fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
	return iss
}

func init() {
	client = resty.New()
	requestCount = 0
	client.SetTimeout(140 * time.Second)
	client.SetHeader("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
	client.SetHeader("Accept", "application/json, text/javascript, */*; q=0.01")
	client.SetHeader("Accept-Language", "zh-CN,zh;q=0.9")
	client.SetHeader("Accept-Encoding", "gzip, deflate")
	client.SetHeader("Host", "macaoapply.singlewindow.gd.cn")
	client.SetHeader("Origin", "https://macaoapply.singlewindow.gd.cn")
	client.SetHeader("Referer", "https://macaoapply.singlewindow.gd.cn/")
	client.SetBaseURL("https://macaoapply.singlewindow.gd.cn")
	// client.SetDebug(true)
	LoadCookie()
}

func GetClient() *resty.Client {
	return client
}

func SaveCookie() {
	var cookies []http.Cookie
	userConf := &config.Config.UserOption
	for _, cookie := range client.Cookies {
		// log.Println("cookie: ", cookie.String())
		cookies = append(cookies, *cookie)
	}

	userConf.Cookies = cookies
	// log.Panicln("cookies: ", userConf.Cookies)
	config.SaveConfig()
}

func SetToken(token string) {
	client.SetHeader("X-Access-Token", token)
}

func LoadCookie() {
	// 先清理
	client.Cookies = make([]*http.Cookie, 0)
	userConf := &config.Config.UserOption
	for _, cookie := range userConf.Cookies {
		client.SetCookie(&cookie)
		// token
		if cookie.Name == "token" {
			SetToken(cookie.Value)
		}
	}
}

type Response struct {
	ResponseCode    int                    `json:"responseCode"`
	ResponseMessage string                 `json:"responseMessage"`
	ResponseResult  map[string]interface{} `json:"responseResult"`
}

func SigningMethodSHA256() jwt.SigningMethod {
	rewrite := jwt.SigningMethodHS256
	rewrite.Name = "SHA256"
	return rewrite
}

func getJwtToken(data jwt.MapClaims) string {
	secret := "kIK0E3eP8GzOGoHrErZQ1BNmMCAwEAAQ==abc"
	header := base64.StdEncoding.EncodeToString([]byte(`{"typ":"JWT","alg":"SHA256"}`))
	payloadJson, _ := json.Marshal(data)
	payload := base64.StdEncoding.EncodeToString(payloadJson)
	hmacInstance := hmac.New(sha256.New, []byte(secret)) // empty secret
	hmacInstance.Write([]byte(header + "." + payload))
	signature := base64.StdEncoding.EncodeToString(hmacInstance.Sum(nil))
	jwtToken := header + "." + payload + "." + signature
	return jwtToken
}

func Request(method string, url string, data jwt.MapClaims) (string, error) {
	if data == nil {
		data = jwt.MapClaims{}
	}
	// 加入iss
	iss := config.Config.UserOption.Iss
	if iss == "" {
		iss = genIss()
		config.Config.UserOption.Iss = iss
		config.SaveConfig()
	}
	data["iss"] = iss
	data["issType"] = "web"
	data["appType"] = "web"

	jwtStr := getJwtToken(data)
	// log.Println("jwtStr: " + jwtStr)
	// 请求开始时间
	startTime := time.Now()

	var err error
	var resp *resty.Response
	if method == "GET" {

		resp, err = client.R().Get(url + "?jwt=" + jwtStr)
		if err != nil {
			return "", err
		}
	} else if method == "POST" {
		req := client.R()
		resp, err = req.SetFormData(map[string]string{
			"jwt": jwtStr,
		}).
			Post(url)
	}
	if err != nil {
		return "", err
	}
	client.Cookies = make([]*http.Cookie, 0)
	client.SetCookies(resp.Cookies())
	code := gjson.GetBytes(resp.Body(), "responseCode").Int()
	msg := gjson.GetBytes(resp.Body(), "responseMessage").String()
	if code != 200 {
		// 错误处理
		// 80x 重新登录
		if code == 801 || code == 802 || code == 803 {
			log.Println("登录过期，重新登录")
			// 重新登录
			Login()
			// 重新请求
			return Request(method, url, data)
		}
		// 70x 验证码
		if code == 701 || code == 702 || code == 703 {
			log.Println("验证码错误，清除验证码缓存")
			// 清除验证码缓存
			cache.ClearCaptchaCache()
		}
		return "", fmt.Errorf("请求失败: %s", msg)
	}
	requestCount++
	// 请求结束时间
	endTime := time.Now()
	// 请求耗时
	costTime := endTime.Sub(startTime)
	log.Println("请求耗时: ", costTime)
	// 请求耗时超过 10s 认为是慢网络
	if costTime > 10*time.Second {
		isSlowNetwork = true
	} else {
		isSlowNetwork = false
	}

	return gjson.GetBytes(resp.Body(), "responseResult").Raw, nil
}

const TimeOutErrText = "Client.Timeout exceeded while awaiting headers"

func RequestWithRetry(method string, url string, data jwt.MapClaims) (string, error) {
	var err error
	var resp string
	for i := 0; i < 5; i++ {
		resp, err = Request(method, url, data)
		if err != nil {
			// 如果不是超时错误，直接返回
			errText := err.Error()
			if !strings.Contains(errText, TimeOutErrText) {
				return "", err
			}
			log.Println("请求超时... 2s后重试")
			// 重试 1s
			time.Sleep(2 * time.Second)
			continue
		}
		return resp, nil
	}
	return "", err
}

func RequestWithCache(method string, url string, data jwt.MapClaims) (string, error) {
	if cache.RequestCache == nil {
		cache.RequestCache = make(map[string]string)
	}
	if cache.RequestCache[url] != "" {
		log.Println("缓存命中")
		return cache.RequestCache[url], nil
	}
	resp, err := RequestWithRetry(method, url, data)
	if err != nil {
		return "", err
	}
	cache.RequestCache[url] = resp
	return resp, nil
}

// 针对慢网络-多线程请求
func RequestWithMulti(method string, url string, data jwt.MapClaims) (string, error) {
	const threadCount = 8
	var resp string
	var err error
	var ch = make(chan string, threadCount)
	for i := 0; i < threadCount; i++ {
		// map不是线程安全的，所以这里需要复制一份
		dataCopy := make(jwt.MapClaims)
		for k, v := range data {
			dataCopy[k] = v
		}
		go func(data jwt.MapClaims) {
			resp, err = RequestWithRetry(method, url, data)
			if err != nil {
				log.Println("请求失败: ", err)
			}
			ch <- resp
		}(dataCopy)
	}
	for i := 0; i < threadCount; i++ {
		resp = <-ch
		if resp != "" {
			break
		}
	}
	return resp, err
}

// 自动选择
func RequestAuto(method string, url string, data jwt.MapClaims) (string, error) {
	if isSlowNetwork {
		log.Println("检测到慢网络，使用多线程请求")
		return RequestWithMulti(method, url, data)
	}
	return RequestWithRetry(method, url, data)
}

func GetToken() *http.Cookie {
	for _, cookie := range client.Cookies {
		if cookie.Name == "token" {
			return cookie
		}
	}
	return nil
}

func IsLogin() bool {
	// token
	token := GetToken()
	if token == nil {
		return false
	}
	return true
}
