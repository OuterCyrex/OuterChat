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
	r.GET("/user/loginByName", service.LoginByName)
	r.GET("/user/getUser", service.GetUser)
	r.POST("/user/add", service.CreateUser)
	Auth := r.Group("/")
	Auth.Use(middleware.Auth())
	{
		Auth.GET("/user/list", service.GetUserList)
		Auth.DELETE("/user/delete", service.DeleteUser)
		Auth.PUT("/user/update", service.UpdateUser)
		Auth.GET("/user/getUserByToken", service.GetUserByToken)

		Auth.GET("/user/sendUserMsg", service.SendUserMsg)
		Auth.GET("/user/history", service.GetHistory)

		Auth.GET("/user/getFriendList", service.GetFriendListById)
		Auth.POST("/user/pushFriendRequest", service.PushFriendRequest)
		Auth.PUT("/user/dealWithFriendRequest", service.DealWithFriendRequest)
		Auth.GET("/user/getRequestWithOption", service.GetRequestWithOption)
		Auth.DELETE("/user/deleteFriend", service.DeleteFriend)
	}
	gin.SetMode(gin.DebugMode)
	return r
}
