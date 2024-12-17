package router

import (
	_ "OuterChat/docs"
	"OuterChat/middleware"
	"OuterChat/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/index", service.GetIndex)
	r.GET("/user/list", service.GetUserList)
	r.POST("/user/add", service.CreateUser)
	r.DELETE("/user/delete", service.DeleteUser)
	r.PUT("/user/update", service.UpdateUser)
	r.GET("/user/getUserByToken", service.GetUserByToken)
	r.GET("/user/loginByName", service.LoginByName)
	r.GET("user/getUser", service.GetUser)

	r.GET("/user/sendMsg", service.SendMsg)
	r.GET("/user/sendUserMsg", service.SendUserMsg)

	r.GET("/user/getFriendList", service.GetFriendListById)
	r.POST("/user/pushFriendRequest", service.PushFriendRequest)
	r.PUT("/user/dealWithFriendRequest", service.DealWithFriendRequest)
	r.GET("/user/getRequestWithOption", service.GetRequestWithOption)
	gin.SetMode(gin.DebugMode)
	return r
}
