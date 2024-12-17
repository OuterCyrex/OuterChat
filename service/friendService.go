package service

import (
	"OuterChat/model"
	"OuterChat/util"
	"OuterChat/util/SError"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetFriendListById
// @Tags 好友
// @Summary 获取用户的好友列表
// @Param id query int false "ID"
// @Success 200 {object} util.Response
// @Router /user/getFriendList [get]
func GetFriendListById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	Friends, err := model.GetFriendListById(uint(id))
	if err != nil {
		util.SendErrorResponse(c, SError.IntervalError, fmt.Sprintf("数据库查询出错：%v", err))
		return
	}
	util.SendSuccessResponse(c, Friends)
}

// PushFriendRequest
// @Tags 好友
// @Summary 发送好友请求
// @Param FromId formData int false "发送者ID"
// @Param TargetId formData int false "接收者ID"
// @Param Desc formData string false "描述"
// @Success 200 {object} util.Response
// @Router /user/pushFriendRequest [post]
func PushFriendRequest(c *gin.Context) {
	FromId, _ := strconv.Atoi(c.PostForm("FromId"))
	TargetId, _ := strconv.Atoi(c.PostForm("TargetId"))
	desc := c.PostForm("Desc")
	if !(model.CheckIdExist(FromId) && model.CheckIdExist(TargetId)) {
		util.SendErrorResponse(c, SError.InValidIdError, "无效 ID")
		return
	}
	contact, err := model.PushFriendRequest(uint(FromId), uint(TargetId), desc)
	if err != nil {
		util.SendErrorResponse(c, SError.IntervalError, err.Error())
		return
	}
	util.SendSuccessResponse(c, contact)
}

// DealWithFriendRequest
// @Tags 好友
// @Summary 处理好友请求, status = 1 为'接受',status = 2 为 '拒绝'
// @Param RequestId query int false "申请ID"
// @Param Status formData int false "接受或拒绝"
// @Success 200 {object} util.Response
// @Router /user/dealWithFriendRequest [put]
func DealWithFriendRequest(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("RequestId"))
	status, _ := strconv.Atoi(c.PostForm("Status"))
	contact, err := model.DealWithFriendRequest(uint(id), status)
	if err != nil {
		util.SendErrorResponse(c, SError.IntervalError, err.Error())
		return
	}
	util.SendSuccessResponse(c, contact)
}

// GetRequestWithOption
// @Tags 好友
// @Summary 获取好友请求, option = 1 为'收到的请求',option = 2 为 '发送的请求'
// @Param Id query int false "用户ID"
// @Param Option query int false "设置"
// @Success 200 {object} util.Response
// @Router /user/getRequestWithOption [get]
func GetRequestWithOption(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("Id"))
	option, _ := strconv.Atoi(c.Query("Option"))
	Result, err := model.GetRequest(uint(id), option)
	if err != nil {
		util.SendErrorResponse(c, SError.IntervalError, err.Error())
		return
	}
	util.SendSuccessResponse(c, Result)
}

// DeleteFriendRequest
// @Tags 好友
// @Summary 删除好友请求
// @
