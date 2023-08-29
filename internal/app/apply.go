package app

import (
	"log"
	"macaoapply-auto/internal/client"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tidwall/gjson"
)

type AppointmentDate struct {
	Quota           int64 `json:"quota"`
	IsFull          bool  `json:"isFull"`
	ApplyNum        int64 `json:"applyNum"`
	AppointmentDate int64 `json:"appointmentDate"`
}

func GetAppointmentDateList() ([]AppointmentDate, error) {
	resp, err := client.Request("POST", "before/sys/appointment/getAppointmentDate", jwt.MapClaims{
		"appointmentType": "passBooking",
		"direction":       "S",
	})
	if err != nil {
		log.Println("获取预约日期失败：" + err.Error())
		return nil, err
	}
	var appointmentDateList []AppointmentDate
	gjson.Get(resp, "appointmentDateList").ForEach(func(_, value gjson.Result) bool {
		appointmentDate := AppointmentDate{
			Quota:           (value.Get("quota").Int()),
			IsFull:          value.Get("isFull").Bool(),
			ApplyNum:        (value.Get("applyNum").Int()),
			AppointmentDate: value.Get("appointmentDate").Int(),
		}
		appointmentDateList = append(appointmentDateList, appointmentDate)
		return true
	})
	return appointmentDateList, nil
}

func CheckAppointmentListHasAvailable(list []AppointmentDate, date string) bool {
	for _, appointmentDate := range list {
		if appointmentDate.Quota > 0 && !appointmentDate.IsFull {
			// 判断是否是指定日期
			if time.Unix(appointmentDate.AppointmentDate, 0).Format("2006-01-02") == date {
				return true
			}
		}
	}
	return false
}

type AppointmentInfo struct {
	PlateNumber     string
	appointmentDate int64
}

// 执行预约
// func DoAppointment(data AppointmentInfo) error {
// 	// 1. 获取 form instance id
// 	formInstance, err := getPassQualification(data.PlateNumber)
// 	if err != nil {
// 		return err
// 	}
// 	// 2. 处理验证码
// 	err = handleCaptcha(formInstance)

// }
