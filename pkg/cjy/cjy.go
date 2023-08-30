package cjy

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"macaoapply-auto/pkg/config"

	"github.com/go-resty/resty/v2"
)

// 超级鹰SDK

var cjyConf *config.CJYOption
var client *resty.Client

type CJYResp struct {
	ErrNo  int    `json:"err_no"`
	ErrStr string `json:"err_str"`
	PicId  string `json:"pic_id"`
	PicStr string `json:"pic_str"`
	Md5    string `json:"md5"`
}

func init() {
	cjyConf = &config.Config.CJYOption
	client = resty.New()
}

func GetPicVal(data []byte) *CJYResp {
	return GetPicRes(data, "1004")
}

func GetPicPos(data []byte) *CJYResp {
	return GetPicRes(data, "9101")
}

func GetPicRes(data []byte, codeType string) *CJYResp {
	url := "https://upload.chaojiying.net/Upload/Processing.php"
	var dataBase64 string
	dataBase64 = base64.StdEncoding.EncodeToString(data)
	resp, err := client.R().
		SetFormData(map[string]string{
			"user":        cjyConf.Username,
			"pass2":       cjyConf.Password,
			"softid":      cjyConf.SoftId,
			"codetype":    codeType,
			"file_base64": dataBase64,
		}).
		Post(url)
	if err != nil {
		log.Println("超级鹰请求出错" + err.Error())
		return nil
	}
	var cjyResp CJYResp
	err = json.Unmarshal(resp.Body(), &cjyResp)
	if err != nil {
		log.Println("超级鹰解析出错" + err.Error())
		return nil
	}
	if cjyResp.ErrNo != 0 {
		log.Println("超级鹰返回错误" + cjyResp.ErrStr)
		return nil
	}
	return &cjyResp
}

func ReportError(picId string) {
	url := "https://upload.chaojiying.net/Upload/ReportError.php"
	_, err := client.R().
		SetFormData(map[string]string{
			"user":   cjyConf.Username,
			"pass2":  cjyConf.Password,
			"softid": cjyConf.SoftId,
			"id":     picId,
		}).
		Post(url)
	if err != nil {
		log.Println("超级鹰报错出错" + err.Error())
	}
}
