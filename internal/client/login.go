package client

import (
	"encoding/base64"
	"log"
	"macaoapply-auto/pkg/cjy"
	"macaoapply-auto/pkg/config"
	"macaoapply-auto/pkg/yunma"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tidwall/gjson"
)

var userConf *config.UserOption

func GetLoginPicVal(data string) (string, string) {
	if config.Config.CaptchaEngine == config.CaptchaEngineCJY {
		imageData, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			log.Println(err)
			return "", ""
		}
		cjyResp := cjy.GetPicVal(imageData)
		if cjyResp == nil {
			log.Println("验证码识别失败")
			return "", ""
		}
		log.Println("验证码识别成功")
		log.Println("验证码：" + cjyResp.PicStr)
		return cjyResp.PicStr, cjyResp.PicId
	} else if config.Config.CaptchaEngine == config.CaptchaEngineYunMa {
		yunmaConf := &config.Config.YunMaOption
		resp, err := yunma.CommonVerify(data, yunmaConf.Token)
		if err != nil {
			log.Println("验证码识别失败" + err.Error())
			return "", ""
		}
		log.Println("resp:" + resp)
		code := gjson.Get(resp, "data").String()
		unId := gjson.Get(resp, "unique_code").String()
		return code, unId
	} else {
		log.Println("未知验证码引擎")
		return "", ""
	}
}

func ReportError(id string) {
	if config.Config.CaptchaEngine == config.CaptchaEngineCJY {
		cjy.ReportError(id)
	} else if config.Config.CaptchaEngine == config.CaptchaEngineYunMa {
		yunmaConf := &config.Config.YunMaOption
		yunma.ReportError(id, yunmaConf.Token)
	}
}

func Login() {
	userConf = &config.Config.UserOption
	// 获取登录验证码
	verifyCode, err := getLoginVerifyCode()
	if err != nil {
		return
	}
	log.Println("登录验证码获取成功")
	// log.Println("验证码地址：" + verifyCode.imageUrl)
	log.Println("验证码ID：" + verifyCode.verifyCodeId)
	// 去除data:image/png;base64,
	imageData := strings.Split(verifyCode.imageUrl, ",")[1]

	code, _ := GetLoginPicVal(imageData)
	if code == "" {
		return
	}
	log.Println("验证码：" + code)

	// 登录
	resp, err := RequestAuto("POST", "before/login", jwt.MapClaims{
		"accountNo":             userConf.Username,
		"password":              userConf.Password,
		"verificationCode":      code,
		"pVerificationCode":     "",
		"loginVerifyCode":       code,
		"verifyCodeId":          verifyCode.verifyCodeId,
		"isNeedCheckVerifyCode": "true",
		"accountType":           "personal",
	})
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "驗證碼錯誤") {
			// 平台报错
			log.Println("验证码错误，报告错误")
			// ReportError(id)
		}
		log.Println("登录失败：" + err.Error())
		return
	}
	token := gjson.Get(resp, "token").String()
	SetToken(token)
	log.Println("登录成功", token)

	SaveCookie() // 保存cookie
}

type getLoginVerifyCodeResp struct {
	imageUrl     string
	verifyCodeId string
}

func getLoginVerifyCode() (getLoginVerifyCodeResp, error) {
	resp, err := Request("GET", "before/sys/verifyCode/getLoginVerifyCode", nil)
	if err != nil {
		log.Println("获取登录验证码失败：" + err.Error())

		return getLoginVerifyCodeResp{}, err
	}
	imageUrl := gjson.Get(resp, "imageUrl").String()
	verifyCodeId := gjson.Get(resp, "verifyCodeId").String()
	return getLoginVerifyCodeResp{
		imageUrl:     imageUrl,
		verifyCodeId: verifyCodeId,
	}, nil
}
