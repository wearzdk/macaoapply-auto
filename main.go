package main

import (
	"log"
	"macaoapply-auto/internal/model"
	"macaoapply-auto/internal/router"
)

func init() {
	//gin.SetMode(gin.ReleaseMode)
}

func main() {
	model.Setup()
	server := router.InitRouter()
	log.Println("server run at 8899")
	err := server.Run(":8899")
	if err != nil {
		return
	}
}
