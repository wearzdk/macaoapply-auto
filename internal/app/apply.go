package app

import (
	"log"
	"macaoapply-auto/internal/cache"
	"macaoapply-auto/internal/client"
	"macaoapply-auto/pkg/config"
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
	resp, err := client.RequestAuto("POST", "before/sys/appointment/getAppointmentDate", jwt.MapClaims{
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
	log.Println("获取预约日期成功")
	// log.Println(appointmentDateList)
	return appointmentDateList, nil
}

func CheckAppointmentListHasAvailable(list []AppointmentDate, date string) bool {
	for _, appointmentDate := range list {
		if !appointmentDate.IsFull {
			// 判断是否是指定日期
			if time.Unix(appointmentDate.AppointmentDate, 0).Format("2006-01-02") == date {
				return true
			}
		}
	}
	return false
}

func DoAppointment(data config.AppointmentOption, formInstance FormInstance) error {
	var err error
	// 2. 处理验证码
	captchaData := handelCaptcha(formInstance.FormInstanceID)
	// 3. validationPassBooking
	err = validationPassBooking(&formInstance, captchaData, &data)
	if err != nil {
		return err
	}
	// 4. doAppointment
	err = createPassAppointment(&formInstance, captchaData, &data)
	if err != nil {
		return err
	}
	return nil
}

func validationPassBooking(formInstance *FormInstance, captchaData cache.CaptchaData, appointmentInfo *config.AppointmentOption) error {
	resp, err := client.RequestWithRetry("POST", "before/sys/appointment/validationPassBooking", jwt.MapClaims{
		"formInstanceId":      formInstance.FormInstanceID,
		"appointmentType":     "passBooking",
		"direction":           "S",
		"plateNumber":         appointmentInfo.PlateNumber,
		"appointmentDate":     appointmentInfo.AppointmentDate,
		"verifyUploadData":    captchaData["verifyUploadData"],
		"checkCaptchaBoolean": true,
		"thisCheckCaptchaId":  captchaData["id"],
	})
	if err != nil {
		log.Println("验证预约失败：" + err.Error())
		return err
	}
	log.Println(resp)
	return nil
}

func createPassAppointment(formInstance *FormInstance, captchaData cache.CaptchaData, appointmentInfo *config.AppointmentOption) error {
	resp, err := client.RequestWithRetry("POST", "before/sys/appointment/createPassAppointment", jwt.MapClaims{
		"formInstanceId":      formInstance.FormInstanceID,
		"appointmentType":     "passBooking",
		"direction":           "S",
		"plateNumber":         appointmentInfo.PlateNumber,
		"appointmentDate":     appointmentInfo.AppointmentDate,
		"verifyUploadData":    captchaData["verifyUploadData"],
		"checkCaptchaBoolean": true,
		"thisCheckCaptchaId":  captchaData["id"],
	})
	if err != nil {
		log.Println("预约失败：" + err.Error())
		return err
	}
	log.Println("预约成功！服务器返回数据：", resp)
	// log.Println(resp)
	return nil
}
