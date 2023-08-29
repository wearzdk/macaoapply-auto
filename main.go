package main

import (
	"gin-mini-starter/internal/model"
	"gin-mini-starter/internal/router"
	"log"
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
