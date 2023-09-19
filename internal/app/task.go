package app

import (
	"context"
	"log"
	"macaoapply-auto/internal/client"
	"macaoapply-auto/pkg/config"
	"strings"
	"time"
)

type TaskFunc func(context.Context) bool

type Task struct {
	fn     TaskFunc
	name   string
	next   *Task
	before *Task
}

// 链表
var taskList *Task

var currentTask *Task

func AddTask(fn TaskFunc, name string) {
	task := &Task{
		fn:   fn,
		name: name,
	}
	if taskList == nil {
		taskList = task
		return
	}
	lastTask := taskList
	for lastTask.next != nil {
		lastTask = lastTask.next
	}
	lastTask.next = task
}

func CheckoutTask(name string) {
	task := taskList
	for task != nil {
		if task.name == name {
			currentTask = task
		}
		task = task.next
	}
}

func NextTask() {
	if currentTask == nil {
		return
	}
	currentTask = currentTask.next
}

func ClearTask() {
	taskList = nil
	currentTask = nil
}

func loginTask(ctx context.Context) bool {
	// 检查是否登录
	if client.IsLogin() {
		return true
	}
	log.Println("未登录，正在登录...")
	Login()
	if client.IsLogin() {
		return true
	}
	log.Println("登录失败")
	ShortWait()
	return false
}

func onLogOut() {
	log.Println("登录过期，重新登录...")
	CheckoutTask("login")
}

var formInstanceGlobe *FormInstance

func getPassQualificationTask(ctx context.Context) bool {
	applyInfo := config.Config.AppointmentOption
	log.Println("正在获取预约资格...")
	var err error
	formInstance, err := getPassQualification(applyInfo.PlateNumber)
	if err != nil {
		log.Println("获取预约资格失败：" + err.Error())
		ShortWait()
		return false
	}
	log.Println("获取预约资格成功" + formInstance.FormInstanceID)
	formInstanceGlobe = &formInstance
	return true
}

func getAppointmentDateListTask(ctx context.Context) bool {
	list, err := GetAppointmentDateList()
	if err != nil {
		log.Println("获取预约日期失败：" + err.Error())
		ShortWait()
		return false
	}
	actDate := time.Unix(config.Config.AppointmentOption.AppointmentDate, 0).Format("2006-01-02")
	if !CheckAppointmentListHasAvailable(list, actDate) {
		log.Println("无可用预约")
		ShortWait()
		return false
	}
	log.Println("有可用预约，正在预约...")
	return true
}

func doAppointmentTask(ctx context.Context) bool {
	applyInfo := config.Config.AppointmentOption
	formInstance := formInstanceGlobe
	if formInstance == nil {
		log.Println("未找到formInstance 回到预约前")
		CheckoutTask("getPassQualification")
		return false
	}
	// 预约
	err := DoAppointment(applyInfo, formInstance)
	if err != nil {
		log.Println("预约失败：" + err.Error())
		errText := err.Error()
		if strings.Contains(errText, "預約名額已滿") {
			log.Println("预约名额已满，回到预约前")
			CheckoutTask("getPassQualification")
			return false
		}
		log.Println("等待30s...")
		time.Sleep(30 * time.Second)
		return false
	}
	log.Println("预约成功！预约进程即将退出...")
	return true
}
