package app

import (
	"encoding/base64"
	"log"
	"macaoapply-auto/internal/client"
	"macaoapply-auto/pkg/cjy"
	"macaoapply-auto/pkg/imageText"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type ComplexImageReturn struct {
	imageData   string
	id          string
	originWidth int
}

func getPassBookingVerifyComplexImage() (ComplexImageReturn, error) {
	resp, err := client.Request("GET", "before/sys/captcha/getPassBookingVerifyComplexImage", nil)
	if err != nil {
		log.Println("获取滑动验证码失败：" + err.Error())
		return ComplexImageReturn{}, err
	}
	log.Println(resp)
	imgUrl := gjson.Get(resp, "responseList.captcha.backgroundImage").String()
	id := gjson.Get(resp, "responseList.id").String()
	originWidth := gjson.Get(resp, "responseList.captcha.backgroundImageWidth").Int()
	// 去除data:image/jpeg;base64,
	return ComplexImageReturn{
		imageData:   imgUrl,
		id:          id,
		originWidth: int(originWidth),
	}, nil
}

// 处理验证码
func handleCaptcha() (map[string]interface{}, error) {
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
	// 识别验证码
	cjyResp := cjy.GetPicVal(newImage)
	if cjyResp == nil {
		return nil, err
	}
	log.Println("验证码识别成功")
	log.Println("中间位置" + cjyResp.PicStr) // x,y
	x, err := strconv.Atoi(strings.Split(cjyResp.PicStr, ",")[0])
	originWidth := resp.originWidth
	// 相对位置
	x = x * 260 / originWidth
	// 相对左上角
	x = x - 24
	if err != nil {
		log.Println("验证码识别失败：" + err.Error())
		return nil, err
	}
	// 模拟滑动
	trackList := GenerateTrack(x)
	startSlidingTime := time.Now().Add(time.Duration(-6) * time.Second)
	lastTrack := trackList[len(trackList)-1]
	firstTrack := trackList[0]
	endSlidingTime := startSlidingTime.Add(time.Duration(lastTrack.T-firstTrack.T) * time.Millisecond)
	verifyUploadData := map[string]interface{}{
		"bgImageWidth":     260,
		"bgImageHeight":    0,
		"startSlidingTime": startSlidingTime.Format("2006-01-02T15:04:05.000Z"),
		"entSlidingTime":   endSlidingTime.Format("2006-01-02T15:04:05.000Z"),
		"trackList":        trackList,
	}
	log.Printf("verifyUploadData: %v\n", verifyUploadData)

	return verifyUploadData, nil

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

	for currentOffset <= target {
		speed := rand.Float64()*2 + 2
		move := rand.Intn(2)
		if currentOffset+move > target {
			move = target - currentOffset
		}
		currentOffset += move
		startTime += int64(speed)

		trackList = append(trackList, Track{
			X:    currentOffset,
			Y:    rand.Intn(2),
			Type: "move",
			T:    startTime,
		})
	}

	trackList = append(trackList, Track{
		X:    currentOffset,
		Y:    0,
		Type: "up",
		T:    startTime + 2,
	})

	return trackList
}
