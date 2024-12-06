package main

import (
	"OuterChat/config"
	"OuterChat/model"
	"OuterChat/router"
)

// @title OuterChat API文档
// @version 1.0
// @description 测试
// @host localhost:8080
// @BasePath /

func main() {
	config.InitConfig()
	model.InitDatabase()
	model.InitCache()
	r := router.Router()
	_ = r.Run(":8080")
}
