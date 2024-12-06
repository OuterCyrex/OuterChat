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
