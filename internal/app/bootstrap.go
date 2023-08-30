package app

import (
	"log"
	"macaoapply-auto/internal/client"
	"math/rand"
	"time"
)

func Wait() {
	log.Println("等待随机6-10秒...")
	sec := rand.Intn(5) + 6
	time.Sleep(time.Duration(sec) * time.Second)
}

func BootStrap() {
	actDate := "2023-09-02"
	actDateUnix, _ := time.Parse("2006-01-02", actDate)
	log.Println("启动...")
	for {
		// 检查是否登录
		if client.IsLogin() {
			break
		}
		log.Println("未登录，正在登录...")
		client.Login()
		if client.IsLogin() {
			break
		}
		log.Println("登录失败")
		Wait()
	}
	log.Println("当前已登录")
	formInstance, err := getPassQualification("MQ-32-83")
	if err != nil {
		log.Println("获取预约资格失败：" + err.Error())
		return
	}
	log.Println("获取预约资格成功" + formInstance.FormInstanceID)
	for {
		list, err := GetAppointmentDateList()
		if err != nil {
			log.Println("获取预约日期失败：" + err.Error())
			Wait()
			continue
		}
		if !CheckAppointmentListHasAvailable(list, actDate) {
			log.Println("无可用预约")
			Wait()
			continue
		}
		log.Println("有可用预约，正在预约...")
		// 预约
		for {
			err = DoAppointment(AppointmentInfo{
				PlateNumber:     "MQ-32-83",
				appointmentDate: actDateUnix.Unix(),
			})
			if err != nil {
				log.Println("预约失败：" + err.Error())
				Wait()
				continue
			}
			log.Println("预约成功")
		}
	}
}
