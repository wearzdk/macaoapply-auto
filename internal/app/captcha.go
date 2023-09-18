package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"macaoapply-auto/internal/cache"
	"macaoapply-auto/internal/client"
	"macaoapply-auto/pkg/cjy"
	"macaoapply-auto/pkg/config"
	"macaoapply-auto/pkg/imageText"
	"macaoapply-auto/pkg/yunma"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tidwall/gjson"
)

type ComplexImageReturn struct {
	imageData   string
	slideImage  string
	id          string
	originWidth int
}

func getPassBookingVerifyComplexImage() (ComplexImageReturn, error) {
	resp, err := client.RequestAuto("GET", "before/sys/captcha/getPassBookingVerifyComplexImage", nil)
	if err != nil {
		log.Println("获取滑动验证码失败：" + err.Error())
		return ComplexImageReturn{}, err
	}
	// log.Println(resp)
	imgUrl := gjson.Get(resp, "responseList.captcha.backgroundImage").String()
	slideImage := gjson.Get(resp, "responseList.captcha.templateImage").String()
	id := gjson.Get(resp, "responseList.id").String()
	originWidth := gjson.Get(resp, "responseList.captcha.backgroundImageWidth").Int()
	// 去除data:image/jpeg;base64,
	return ComplexImageReturn{
		imageData:   imgUrl[strings.Index(imgUrl, ",")+1:],
		slideImage:  slideImage[strings.Index(slideImage, ",")+1:],
		id:          id,
		originWidth: int(originWidth),
	}, nil
}

// 验证
func checkPassBookingComplexImage(id string, formInstanceId string, data cache.CaptchaData) (bool, error) {
	_, err := client.RequestAuto("POST", "before/sys/captcha/checkPassBookingComplexImage", jwt.MapClaims{
		"formInstanceId":   formInstanceId,
		"appointmentType":  "passBooking",
		"direction":        "S",
		"verifyUploadData": data,
		"id":               id,
	})
	if err != nil {
		log.Println("验证滑动验证码失败：" + err.Error())
		return false, err
	}
	// log.Println(resp)
	return true, nil
}

type ProcessCaptchaFunc func(resp *ComplexImageReturn) (x float64, err error)

func processCaptchaCjy(data *ComplexImageReturn) (x float64, err error) {
	imageBytes, err := base64.StdEncoding.DecodeString(data.imageData)
	if err != nil {
		log.Println("base64解码失败：" + err.Error())
		return 0, err
	}
	// 插入文字
	newImage, err := imageText.InsertTextToImage(imageBytes, "请点击凹槽正中间")
	if err != nil {
		log.Println("插入文字失败：" + err.Error())
		return 0, err
	}
	// 保存图片
	os.WriteFile("captcha.jpg", newImage, 0666)
	log.Println("图片准备完成")
	// time.Sleep(10 * time.Minute)
	// 识别验证码
	cjyResp := cjy.GetPicPos(newImage)
	if cjyResp == nil {
		return 0, err
	}
	log.Println("验证码识别成功")
	log.Println("中间位置" + cjyResp.PicStr) // x,y
	x_str, err := strconv.Atoi(strings.Split(cjyResp.PicStr, ",")[0])
	if err != nil {
		log.Println("验证码识别失败：" + err.Error())
		return 0, err
	}
	x = float64(x_str)
	originWidth := data.originWidth
	// 相对位置
	x = x * 260 / float64(originWidth)
	// 相对左上角
	x = x - 24.25
	log.Println("相对位置" + strconv.FormatFloat(x, 'f', 2, 64))
	return x, nil
}

func processCaptchaYm(data *ComplexImageReturn) (x float64, err error) {
	ymConf := config.Config.YunMaOption
	if ymConf.Token == "" {
		log.Println("请配置云码token")
		return 0, fmt.Errorf("请配置云码token")
	}
	// 识别验证码
	yunmaResp, err := yunma.SlideVerify(data.slideImage, data.imageData, ymConf.Token)
	if err != nil {
		log.Println("云码识别失败：" + err.Error())
		return 0, err
	}
	x_str := gjson.Get(yunmaResp, "data").String()
	x_int, err := strconv.Atoi(x_str)
	if err != nil {
		log.Println("云码识别失败：数据格式不正确" + err.Error())
		return 0, err
	}
	x = float64(x_int)
	originWidth := data.originWidth
	// 相对位置
	x = x * 260 / float64(originWidth)
	// // 相对左上角
	// x = x - 24.25
	log.Println("相对位置" + strconv.FormatFloat(x, 'f', 2, 64))
	return x, nil
}

func doCaptcha(formInstanceId string, processCaptcha ProcessCaptchaFunc) (cache.CaptchaData, error) {
	// 获取滑动验证码
	resp, err := getPassBookingVerifyComplexImage()
	if err != nil {
		log.Println("获取滑动验证码失败：" + err.Error())
		return nil, err
	}
	startSlidingTime := time.Now()
	// 识别验证码
	x, err := processCaptcha(&resp)
	if err != nil {
		log.Println("验证码识别失败：" + err.Error())
		return nil, err
	}
	log.Println("相对位置" + strconv.FormatFloat(x, 'f', 2, 64))
	// 模拟滑动
	trackList := GenerateTrack(int(x))
	lastTrack := trackList[len(trackList)-1]
	firstTrack := trackList[0]
	startSlidingTime = startSlidingTime.Add(time.Duration(firstTrack.T) * time.Millisecond)
	endSlidingTime := startSlidingTime.Add(time.Duration(lastTrack.T-firstTrack.T) * time.Millisecond)
	verifyUploadData := cache.CaptchaData{
		"bgImageWidth":     260,
		"bgImageHeight":    159,
		"startSlidingTime": startSlidingTime.Format("2006-01-02T15:04:05.000Z"),
		"entSlidingTime":   endSlidingTime.Format("2006-01-02T15:04:05.000Z"),
		"trackList":        trackList,
	}
	verifyUploadDataJson, _ := json.Marshal(verifyUploadData)
	log.Println("verifyUploadDataJson: ", string(verifyUploadDataJson))
	// 检查是否到endSlidingTime
	// for {
	// 	if time.Now().After(endSlidingTime) {
	// 		break
	// 	}
	// 	log.Println("等待滑动验证码结束")
	// 	time.Sleep(endSlidingTime.Sub(time.Now()))
	// }
	// 验证
	ok, err := checkPassBookingComplexImage(resp.id, formInstanceId, verifyUploadData)
	if err != nil {
		log.Println("验证滑动验证码失败：" + err.Error())
		return nil, err
	}
	if !ok {
		log.Println("验证滑动验证码失败")
		return nil, err
	}
	log.Println("验证滑动验证码成功")
	return cache.CaptchaData{
		"id":               resp.id,
		"verifyUploadData": verifyUploadData,
	}, nil
}

// 处理验证码
func handelCaptcha(formInstanceId string) cache.CaptchaData {
	// 从缓存中获取
	if cache.CaptchaCache != nil {
		return cache.CaptchaCache
	}
	for {
		var data cache.CaptchaData
		var err error
		if config.Config.CaptchaEngine == config.CaptchaEngineCJY {
			data, err = doCaptcha(formInstanceId, processCaptchaCjy)
		} else if config.Config.CaptchaEngine == config.CaptchaEngineYunMa {
			data, err = doCaptcha(formInstanceId, processCaptchaYm)
		} else {
			log.Println("未知验证码引擎")
			return nil
		}
		if err != nil {
			log.Println("处理验证码失败：1s后重试" + err.Error())
			time.Sleep(1 * time.Second)
			continue
		}
		cache.CaptchaCache = data
		return data
	}
}

type Track struct {
	X    int
	Y    int
	Type string
	T    int64
}

func GenerateTrack(distance int) []Track {
	track := []Track{}
	current := 0
	mid := distance * 4 / 5
	t := 0.2
	v := 0.0
	startTime := int64(2000 + rand.Intn(2000))

	track = append(track, Track{
		X:    0,
		Y:    0,
		Type: "down",
		T:    startTime,
	})

	for {
		a := 0.0
		if current < mid {
			a = 2.0
		} else {
			a = -3.0
		}
		speed := rand.Float64()*5 + 10
		startTime += int64(speed)

		v0 := v
		v = v0 + a*t
		move := v0*t + 0.5*a*t*t
		current += int(move)

		track = append(track, Track{
			X:    current + rand.Intn(2),
			Y:    rand.Intn(2),
			Type: "move",
			T:    startTime,
		})

		if current >= distance {
			break
		}
	}

	track = append(track, Track{
		X:    distance,
		Y:    rand.Intn(2),
		Type: "up",
		T:    startTime + 2,
	})

	return track
}
