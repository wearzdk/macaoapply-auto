package config

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// 配置文件

type CJYOption struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	SoftId   string `json:"softId,omitempty"`
}

type UserOption struct {
	Username string        `json:"username,omitempty"`
	Password string        `json:"password,omitempty"`
	Iss      string        `json:"iss,omitempty"`
	Cookies  []http.Cookie `json:"cookies,omitempty"`
}

// 预约配置
type AppointmentOption struct {
	PlateNumber     string `json:"plateNumber,omitempty"`
	AppointmentDate int64  `json:"appointmentDate,omitempty"`
	StartTime       int64  `json:"startTime,omitempty"`
	EndTime         int64  `json:"endTime,omitempty"`
}

type Option struct {
	CJYOption         `json:"cjy,omitempty"`
	UserOption        `json:"user,omitempty"`
	AppointmentOption `json:"appointment,omitempty"`
	UA                string `json:"ua,omitempty"`
	Port              string `json:"port,omitempty"`
	Speed             int64  `json:"speed,omitempty"`
	Thread            int    `json:"thread,omitempty"`
	OnMulti           bool   `json:"onMulti,omitempty"`
}

var Config Option

func init() {
	// 初始化
	Config = Option{
		Port:   "8080",
		UA:     "Mozilla/5.0 (Linux; Android 10; Redmi K30 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
		Speed:  800,
		Thread: 2,
	}
	// 在配置中读取
	file := FileReading("config.json")
	if file != nil {
		err := json.Unmarshal(file, &Config)
		if err != nil {
			log.Panicf("配置文件解析失败: %s", err.Error())
		}
	}
	// 保存配置
	SaveConfig()
}

func FileSaving(name string, file []byte) {
	// 检查config目录是否存在
	if _, err := os.Stat("config"); os.IsNotExist(err) {
		// 不存在则创建
		err = os.Mkdir("config", 0755)
		if err != nil {
			log.Panicf("创建配置文件目录失败: %s", err.Error())
		}
	}
	// 保存
	err := os.WriteFile("config/"+name, file, 0644)
	if err != nil {
		log.Panicf("保存配置文件失败: %s", err.Error())
	}
}

func FileReading(name string) []byte {
	file, err := os.ReadFile("config/" + name)
	if err != nil {
		return nil
	}
	return file
}

// SaveConfig 保存配置
func SaveConfig() {
	// 序列化
	data, err := json.MarshalIndent(Config, "", "  ")
	if err != nil {
		log.Panicf("配置文件序列化失败: %s", err.Error())
		return
	}
	// 保存
	FileSaving("config.json", data)
}
