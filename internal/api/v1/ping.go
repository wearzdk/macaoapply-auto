package v1

import (
	"log"
	"macaoapply-auto/internal/app"
	"macaoapply-auto/pkg/resp"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// 重启预约
func Restart(c *gin.Context) {
	log.Println("指令已发送，等待进程退出中...请稍后")
	app.Quit()
	time.Sleep(1 * time.Second)
	go app.BootStrap()
	log.Println("已重新启动")
	resp.Success(c)
}
