package app

import (
	"encoding/base64"
	"log"
	"macaoapply-auto/internal/cache"
	"macaoapply-auto/internal/client"
	"macaoapply-auto/pkg/cjy"
	"macaoapply-auto/pkg/imageText"
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
	id := gjson.Get(resp, "responseList.id").String()
	originWidth := gjson.Get(resp, "responseList.captcha.backgroundImageWidth").Int()
	// 去除data:image/jpeg;base64,
	return ComplexImageReturn{
		imageData:   imgUrl[strings.Index(imgUrl, ",")+1:],
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

// 处理验证码
func doCaptcha(formInstanceId string) (cache.CaptchaData, error) {
	// 获取滑动验证码
	resp, err := getPassBookingVerifyComplexImage()
	if err != nil {
		log.Println("获取滑动验证码失败：" + err.Error())
		return nil, err
	}
	imageBytes, err := base64.StdEncoding.DecodeString(resp.imageData)
	if err != nil {
		log.Println("base64解码失败：" + err.Error())
		return nil, err
	}
	// 插入文字
	newImage, err := imageText.InsertTextToImage(imageBytes, "请点击凹槽正中间")
	if err != nil {
		log.Println("插入文字失败：" + err.Error())
		return nil, err
	}
	// 保存图片
	os.WriteFile("captcha.jpg", newImage, 0666)
	log.Println("图片准备完成")
	// time.Sleep(10 * time.Minute)
	// 识别验证码
	cjyResp := cjy.GetPicPos(newImage)
	if cjyResp == nil {
		return nil, err
	}
	log.Println("验证码识别成功")
	log.Println("中间位置" + cjyResp.PicStr) // x,y
	x, err := strconv.Atoi(strings.Split(cjyResp.PicStr, ",")[0])
	if err != nil {
		log.Println("验证码识别失败：" + err.Error())
		return nil, err
	}
	originWidth := resp.originWidth
	// 相对位置
	x = x * 260 / originWidth
	// 相对左上角
	x = x - 25
	log.Println("相对位置" + strconv.Itoa(x))
	// 模拟滑动
	trackList := GenerateTrack(x)
	startSlidingTime := time.Now().Add(time.Duration(-6) * time.Second)
	lastTrack := trackList[len(trackList)-1]
	firstTrack := trackList[0]
	endSlidingTime := startSlidingTime.Add(time.Duration(lastTrack.T-firstTrack.T) * time.Millisecond)
	verifyUploadData := cache.CaptchaData{
		"bgImageWidth":     260,
		"bgImageHeight":    159,
		"startSlidingTime": startSlidingTime.Format("2006-01-02T15:04:05.000Z"),
		"entSlidingTime":   endSlidingTime.Format("2006-01-02T15:04:05.000Z"),
		"trackList":        trackList,
	}
	// log.Printf("verifyUploadData: %v\n", verifyUploadData)
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
		data, err := doCaptcha(formInstanceId)
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

func GenerateTrack(target int) []Track {
	var trackList []Track
	currentOffset := 0
	// 2000 + rand.Intn(1000)
	startTime := int64(2000 + rand.Intn(2000))

	trackList = append(trackList, Track{
		X:    0,
		Y:    0,
		Type: "down",
		T:    startTime,
	})

	for {
		// log.Println("currentOffset: " + strconv.Itoa(currentOffset) + " target: " + strconv.Itoa(target) + " startTime: " + strconv.FormatInt(startTime, 10) + " speed: " + strconv.FormatFloat(rand.Float64()*2+2, 'f', 2, 64))
		speed := rand.Float64()*2 + 2
		move := rand.Intn(2)
		currentOffset += move
		startTime += int64(speed)

		trackList = append(trackList, Track{
			X:    currentOffset,
			Y:    rand.Intn(2),
			Type: "move",
			T:    startTime,
		})
		if currentOffset >= target {
			break
		}
	}

	trackList = append(trackList, Track{
		X:    currentOffset,
		Y:    0,
		Type: "up",
		T:    startTime + 2,
	})

	return trackList
}
