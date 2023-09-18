package yunma

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

const API_URL = "http://api.jfbym.com/api/YmServer/customApi"

func getCodeMsg(code int64) string {
	switch code {
	case 10001:
		return "参数错误"
	case 10002:
		return "余额不足"
	case 10003:
		return "无此访问权限"
	case 10004:
		return "无此验证类型"
	case 10005:
		return "网络拥塞"
	case 10006:
		return "数据包过载"
	case 10007:
		return "服务繁忙"
	case 10008:
		return "网络错误，请稍后重试"
	case 10009:
		return "结果准备中，请稍后再试"
	case 10010:
		return "请求结束"
	default:
		return "未知错误"
	}
}
func postRequest(config map[string]interface{}) (string, error) {
	startTime := time.Now()
	configData, _ := json.Marshal(config)
	body := bytes.NewBuffer([]byte(configData))
	resp, err := http.Post(API_URL, "application/json;charset=utf-8", body)
	if err != nil {
		log.Println("云码请求出错" + err.Error())
		return "", err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	log.Println("云码请求耗时: ", time.Now().Sub(startTime))
	respCode := gjson.Get(string(data), "code").Int()
	if respCode != 10000 {
		respMsg := gjson.Get(string(data), "msg").String()
		praseMsg := getCodeMsg(respCode)
		msg := praseMsg
		if respMsg != "" {
			msg = respMsg
		}
		log.Println("云码返回错误: ", msg)
		return "", fmt.Errorf(msg)
	}
	respData := gjson.Get(string(data), "data").Raw
	return respData, nil
}

func SlideVerify(slideImage string, backgroundImage string, token string) (string, error) {
	config := map[string]interface{}{
		"slide_image":      slideImage,
		"background_image": backgroundImage,
		"type":             "20111",
		"token":            token,
	}
	return postRequest(config)
}

func CommonVerify(image string, token string) (string, error) {
	config := map[string]interface{}{
		"image": image,
		"type":  "10103",
		"token": token,
	}
	return postRequest(config)
}

func ReportError(id string, token string) {
	url := "http://api.jfbym.com/api/YmServer/refundApi"
	config := map[string]interface{}{
		"uniqueCode": id,
		"token":      token,
	}
	configData, _ := json.Marshal(config)
	body := bytes.NewBuffer([]byte(configData))
	_, err := http.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		log.Println("云码报错出错" + err.Error())
	}
}
