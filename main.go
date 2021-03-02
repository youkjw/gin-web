package main

import (
	"fmt"
	"gin-web/pkg/app"
	"gin-web/pkg/setting"
	"gin-web/routers"
	"github.com/fvbock/endless"
	"log"
)

func main() {

	app.Init()

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