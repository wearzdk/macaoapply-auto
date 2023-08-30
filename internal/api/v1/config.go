package v1

import (
	"macaoapply-auto/pkg/config"
	"macaoapply-auto/pkg/resp"

	"github.com/gin-gonic/gin"
)

func GetConfig(c *gin.Context) {
	resp.SuccessData(c, config.Config)
}

func SetConfig(c *gin.Context) {
	if err := c.ShouldBindJSON(&config.Config); err != nil {
		resp.Error(c, resp.CodeInternalServer, err.Error())
		return
	}
	config.SaveConfig()
	resp.Success(c)
}
