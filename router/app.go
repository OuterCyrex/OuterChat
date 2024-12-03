package router

import (
	_ "OuterChat/docs"
	"OuterChat/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/index", service.GetIndex)
	r.GET("/user/list", service.GetUserList)
	r.POST("/user/add", service.CreateUser)
	r.DELETE("/user/delete", service.DeleteUser)
	r.PUT("/user/update", service.UpdateUser)
	r.GET("/user/loginByName", service.LoginByName)

	r.GET("/user/sendMsg", service.SendMsg)
	gin.SetMode(gin.DebugMode)
	return r
}
