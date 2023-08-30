package router

import (
	v1 "macaoapply-auto/internal/api/v1"
	"macaoapply-auto/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// cors
	r.Use(middleware.Cors())

	// Common
	apiV1 := r.Group("/api/v1")
	{
		// ping
		apiV1.GET("/ping", v1.Ping)
		// user
		apiV1.POST("/user/login", v1.LoginUser)
		// config
		apiV1.GET("/config", v1.GetConfig)
		apiV1.POST("/config", v1.SetConfig)

		// restart
		apiV1.POST("/restart", v1.Restart)
	}

	return r
}
