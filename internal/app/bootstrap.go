package app

import (
	"context"
	"log"
	"macaoapply-auto/internal/client"
	"macaoapply-auto/pkg/config"
	"time"
)

var instance context.Context
var cancelFunc context.CancelFunc

func Start() {
	if Running() {
		log.Println("已经在运行中")
		return
	}
	instance, cancelFunc = context.WithCancel(context.Background())
	go BootStrap(instance)
}

func Running() bool {
	return instance != nil
}

func Quit() {
	if cancelFunc != nil {
		cancelFunc()
		cancelFunc = nil
		instance = nil
		log.Println("已退出")
	}
}

func Restart() {
	Quit()
	Start()
}

func ShortWait() {
	log.Printf("等待%d毫秒...\n", config.Config.Speed)
	time.Sleep(time.Duration(config.Config.Speed) * time.Millisecond)
}

func CheckTime() bool {
	startTime := time.Unix(config.Config.AppointmentOption.StartTime, 0)
	// endTime := time.Unix(config.Config.AppointmentOption.EndTime, 0)

	// 如果未到开始时间，等待
	if time.Now().Before(startTime) {
		log.Println("未到开始时间，等待...")
		for {
			startTime = time.Unix(config.Config.AppointmentOption.StartTime, 0)

			if time.Now().After(startTime) {
				break
			}
			// 距离开始时间 取整
			log.Println("距离开始时间还有", int(startTime.Sub(time.Now()).Seconds()), "秒")
			time.Sleep(1 * time.Second)
		}
	}
	return true
}

// 配置检查
func CheckConfig() bool {
	// 超级鹰
	if config.Config.CJYOption.Username == "" || config.Config.CJYOption.Password == "" || config.Config.CJYOption.SoftId == "" {
		log.Println("请先配置超级鹰")
		return false
	}
	// 用户
	if config.Config.UserOption.Username == "" || config.Config.UserOption.Password == "" {
		log.Println("请先配置账户信息")
		return false
	}
	// 预约
	if config.Config.AppointmentOption.PlateNumber == "" || config.Config.AppointmentOption.AppointmentDate == 0 {
		log.Println("请先配置预约信息")
		return false
	}
	return true
}

func BootStrap(ctx context.Context) {
	// 配置检查
	if !CheckConfig() {
		return
	}
	log.Println("启动...")
	if !CheckTime() {
		log.Println("退出...")
		return
	}
	// 清空任务
	ClearTask()
	// 检查是否登录
	AddTask(loginTask, "login")
	// 获取预约资格
	AddTask(getPassQualificationTask, "getPassQualification")
	// 获取预约日期
	AddTask(getAppointmentDateListTask, "getAppointmentDateList")
	// 预约
	AddTask(doAppointmentTask, "doAppointment")
	// 执行任务
	CheckoutTask("login")

	// 配置登出回调
	client.OnLogout = onLogOut

	for {
		if currentTask == nil {
			log.Println("任务执行完毕 预约进程退出...")
			return
		}
		for {
			select {
			case <-ctx.Done():
				return // 如果context已经被取消，则返回
			default:
				ok := currentTask.fn(ctx)
				if ok {
					NextTask()
				}
			}
		}
	}
}
