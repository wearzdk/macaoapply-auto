package v1

import (
	"log"
	"macaoapply-auto/internal/app"
	"macaoapply-auto/pkg/resp"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// 重启
func Restart(c *gin.Context) {
	log.Println("指令已发送，等待进程重启中...请稍后")
	app.Restart()
	resp.Success(c)
}

// 退出预约
func Quit(c *gin.Context) {
	log.Println("指令已发送，等待进程退出中...请稍后")
	app.Quit()
	resp.Success(c)
}

// 启动预约
func Start(c *gin.Context) {
	app.Start()
	resp.Success(c)
}

// 获取预约状态
func Status(c *gin.Context) {
	resp.SuccessData(c, gin.H{
		"running": app.Running(),
	})
}
