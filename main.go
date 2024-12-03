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

// https://www.bilibili.com/video/BV18G411P7Cw?spm_id_from=333.788.videopod.episodes&vd_source=a4e28f82e6d2c94a274a1152798412d4&p=23
