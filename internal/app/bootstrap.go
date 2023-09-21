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

// 配置检查
func CheckConfig() bool {
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

	// 超级鹰
	if config.Config.CaptchaEngine == config.CaptchaEngineCJY && (config.Config.CJYOption.Username == "" || config.Config.CJYOption.Password == "" || config.Config.CJYOption.SoftId == "") {
		log.Println("请先配置超级鹰")
		return false
	}
	// 云码
	if config.Config.CaptchaEngine == config.CaptchaEngineYunMa && (config.Config.YunMaOption.Token == "") {
		log.Println("请先配置云码")
		return false
	}
	return true
}

func BootStrap(ctx context.Context) {
	defer Quit()
	// 配置检查
	if !CheckConfig() {
		return
	}
	log.Println("启动...")
	// 清空任务
	ClearTask()
	// 时间检测
	AddTask(CheckTimeTask, "timeCheck")
	// 检查是否登录
	AddTask(loginTask, "login")
	// 获取预约资格
	AddTask(getPassQualificationTask, "getPassQualification")
	// 测试滑动验证码
	// AddTask(testCaptchaTask, "testCaptcha")
	// 获取预约日期
	AddTask(getAppointmentDateListTask, "getAppointmentDateList")
	// 预约
	AddTask(doAppointmentTask, "doAppointment")
	// 执行任务
	CheckoutTask("timeCheck")

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
				if currentTask == nil {
					log.Println("任务执行完毕 预约进程退出...")
					return
				}
				ok := currentTask.fn(ctx)
				if ok {
					NextTask()
				}
			}
		}
	}
}
