package service

import (
	"OuterChat/model"
	"OuterChat/util"
	"OuterChat/util/SError"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Tags 用户模块
// @Summary 获取用户列表
// @Param Authorization header string false "token"
// @Success 200 {object} util.Response
// @Router /user/list [get]
func GetUserList(c *gin.Context) {
	data, err := model.GetUserList()
	if err != nil {
		util.SendErrorResponse(c, SError.IntervalError, err.Error())
	}
	util.SendSuccessResponse(c, data)
}

// GetUserByToken
// @Tags 用户模块
// @Summary 解析用户token获取信息
// @Param Authorization header string false "token"
// @Success 200 {object} util.Response
// @Router /user/getUserByToken [get]
func GetUserByToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	claim, err := util.ParseToken(tokenString)
	if err != nil {
		util.SendErrorResponse(c, SError.IntervalError, fmt.Sprintf("token解析出错：%v", err))
		return
	}
	user, err := model.FindUserByField("id", claim.UID)
	if err != nil {
		util.SendErrorResponse(c, SError.IntervalError, fmt.Sprintf("数据库查询出错：%v", err))
		return
	}
	util.SendSuccessResponse(c, user)
}

// GetUser
// @Tags 用户模块
// @Summary 通过用户ID获取用户对象
// @Param id query int false "id"
// @Success 200 {object} util.Response
// @Router /user/getUser [get]
func GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	user, err := model.GetUser(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.SendErrorResponse(c, SError.InValidIdError, "无效Id")
			return
		}
		util.SendErrorResponse(c, SError.IntervalError, fmt.Sprintf("数据库查询出错: %v", err))
		return
	}
	util.SendSuccessResponse(c, user)
}

// CreateUser
// @Tags 用户模块
// @Summary 添加用户
// @Param name formData string false "用户名"
// @Param email formData string false "电子邮箱"
// @Param password formData string false "密码"
// @Param repassword formData string false "重复密码"
// @Success 200 {object} util.Response
// @Router /user/add [post]
func CreateUser(c *gin.Context) {

	//password check
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")
	if password != repassword {
		c.JSON(-1, util.ErrorHttpResponse(SError.RePasswordError, "两次密码不一致"))
		c.Abort()
		return
	}

	email := c.PostForm("email")
	if !govalidator.IsEmail(email) || !model.CheckEmailValid(email) {
		c.JSON(-1, util.ErrorHttpResponse(SError.InValidEmailError, "无效邮箱"))
		c.Abort()
		return
	}

	name := c.PostForm("name")
	if !model.CheckNameValid(name) {
		c.JSON(-1, util.ErrorHttpResponse(SError.NameHasBeenUsedError, "用户名已被使用"))
		c.Abort()
		return
	}

	user := model.UserBasic{
		Name:          name,
		Email:         email,
		LoginTime:     time.Now(),
		HeartbeatTime: time.Now(),
		LogoutTime:    time.Now(),
		Password:      util.Md5Encode(password),
	}

	err := model.CreateUser(user).Error
	if err != nil {
		c.JSON(-1, util.ErrorHttpResponse(SError.IntervalError, fmt.Sprintf("数据库错误: %v", err)))
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, util.SuccessHttpResponse(user))
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @Param id query string false "id"
// @Param Authorization header string false "token"
// @Success 200 {object} util.Response
// @Router /user/delete [delete]
func DeleteUser(c *gin.Context) {
	user := model.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	if !model.CheckIdExist(id) {
		c.JSON(-1, util.ErrorHttpResponse(SError.InValidIdError, "无效ID"))
		c.Abort()
		return
	}
	user.ID = uint(id)
	err := model.DeleteUser(user).Error
	if err != nil {
		c.JSON(-1, util.ErrorHttpResponse(SError.IntervalError, fmt.Sprintf("数据库错误: %v", err)))
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, util.SuccessHttpResponse(user))
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @Param id query int false "id"
// @Param name formData string false "用户名"
// @Param password formData string false "密码"
// @Param Authorization header string false "token"
// @Success 200 {object} util.Response
// @Router /user/update [put]
func UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	if !model.CheckIdExist(id) {
		c.JSON(-1, util.ErrorHttpResponse(SError.InValidIdError, "无效ID"))
		c.Abort()
		return
	}

	user := model.UserBasic{
		Name:     c.PostForm("name"),
		Password: util.Md5Encode(c.PostForm("password")),
	}
	user.ID = uint(id)
	err := model.UpdateUser(user).Error
	if err != nil {
		c.JSON(-1, util.ErrorHttpResponse(SError.IntervalError, fmt.Sprintf("数据库错误: %v", err)))
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, util.SuccessHttpResponse(user))
}

// LoginByName
// @Summary 根据用户名登陆
// @Tags 用户模块
// @Param name query string false "用户名"
// @Param password query string false "密码"
// @Success 200 {object} util.Response
// @Router /user/loginByName [get]
func LoginByName(c *gin.Context) {
	name := c.Query("name")
	if model.CheckNameValid(name) {
		util.SendErrorResponse(c, SError.NameNotExistError, "用户名不存在")
		return
	}
	password := c.Query("password")
	err := model.LoginByName(name, util.Md5Encode(password)).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		util.SendErrorResponse(c, SError.PasswordWrongError, "密码错误")
		return
	}

	user, err := model.FindUserByField("Name", name)
	if err != nil {
		util.SendErrorResponse(c, SError.IntervalError, fmt.Sprintf("数据库出错：%v", err))
		return
	}
	token, _ := util.CreateToken(user.ID)
	util.SendSuccessResponse(c, token)
}

func SendUserMsg(c *gin.Context) {
	model.Chat(c.Writer, *c.Request)
}

// GetHistory
// @Summary 获取用户与某人的的聊天记录
// @Tags 用户模块
// @Param FromId query int false "用户ID"
// @Param TargetId query int false "目标ID"
// @Param Authorization header string false "token"
// @Success 200 {object} util.Response
// @Router /user/history [get]
func GetHistory(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("FromId"))
	targetId, _ := strconv.Atoi(c.Query("TargetId"))
	if !(model.CheckIdExist(userId) && model.CheckIdExist(targetId)) {
		util.SendErrorResponse(c, SError.InValidIdError, "无效 ID")
		return
	}
	if !model.IsFriendStatus(uint(userId), uint(targetId), model.WithStatus(model.Accept)) {
		util.SendErrorResponse(c, SError.NotEvenFriendError, "与对方尚且不是好友")
		return
	}
	history, err := model.GetHistory(uint(userId), uint(targetId))
	if err != nil {
		util.SendErrorResponse(c, SError.IntervalError, err.Error())
		return
	}
	util.SendSuccessResponse(c, history)
}
