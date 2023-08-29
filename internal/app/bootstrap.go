package app

import (
	"log"
	"macaoapply-auto/internal/client"
	"time"
)

func BootStrap() {
	actDate := "2023-09-30"
	log.Println("启动...")
	for {
		// 检查是否登录
		if client.IsLogin() {
			break
		}
		log.Println("未登录，正在登录...")
		Login()
		if client.IsLogin() {
			break
		}
		log.Println("登录失败，10s后重试...")
		time.Sleep(10 * time.Second)
	}
	token := client.GetToken()
	log.Println("当前已登录 过期时间：" + token.Expires.Format("2006-01-02 15:04:05"))
	_, err := handleCaptcha()
	if err != nil {
		log.Println("处理验证码失败：" + err.Error())
		return
	}
	log.Println("处理验证码成功")
	for {
		list, err := GetAppointmentDateList()
		if err != nil {
			log.Println("获取预约日期失败，10s后重试：" + err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		if !CheckAppointmentListHasAvailable(list, actDate) {
			log.Println("无可用预约，10s后重试...")
			time.Sleep(10 * time.Second)
			continue
		}
		log.Println("有可用预约，正在预约...")
		// 预约
		for {

		}
	}
}
