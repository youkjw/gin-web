package main

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/logging"
	"gin-web/pkg/redis"
	"gin-web/pkg/setting"
	"gin-web/routers"
	"github.com/fvbock/endless"
	"log"
)

func main() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	redis.Setup()

	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	server := endless.NewServer(endPoint, routers.InitRouter())

	err := server.ListenAndServe()
	if err != nil  {
		log.Printf("Server err: %v", err)
	}
}