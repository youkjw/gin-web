package app

import (
	"gin-web/models"
	"gin-web/pkg/logging"
	"gin-web/pkg/redis"
	"gin-web/pkg/setting"
)

func Init() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	redis.Setup()
}