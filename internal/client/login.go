package client

import (
	"encoding/base64"
	"log"
	"macaoapply-auto/pkg/cjy"
	"macaoapply-auto/pkg/config"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tidwall/gjson"
)

var userConf *config.UserOption

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
	imgData, err := base64.StdEncoding.DecodeString(verifyCode.imageUrl[22:])
	if err != nil {
		log.Println(err)
		return
	}
	// 识别验证码
	cjyResp := cjy.GetPicVal(imgData)
	if cjyResp == nil {
		return
	}
	log.Println("验证码识别成功")
	log.Println("验证码：" + cjyResp.PicStr)

	// 登录
	resp, err := RequestAuto("POST", "before/login", jwt.MapClaims{
		"accountNo":             userConf.Username,
		"password":              userConf.Password,
		"verificationCode":      cjyResp.PicStr,
		"pVerificationCode":     "",
		"loginVerifyCode":       cjyResp.PicStr,
		"verifyCodeId":          verifyCode.verifyCodeId,
		"isNeedCheckVerifyCode": "true",
		"accountType":           "personal",
	})
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "驗證碼錯誤") {
			// 超级鹰报错
			log.Println("验证码错误，汇报超级鹰")
			cjy.ReportError(cjyResp.PicId)
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
