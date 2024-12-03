package service

import (
	"OuterChat/model"
	"OuterChat/util"
	"OuterChat/util/SError"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

var TypeValid = model.UserTypeValidImpl{}

// GetUserList
// @Tags 用户模块
// @Summary 获取用户列表
// @Success 200 {object} util.Response
// @Router /user/list [get]
func GetUserList(c *gin.Context) {
	var data []model.UserBasic
	data = model.GetUserList()
	c.JSON(http.StatusOK, util.SuccessHttpResponse(data))
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
	if !govalidator.IsEmail(email) || !TypeValid.CheckEmailValid(email) {
		c.JSON(-1, util.ErrorHttpResponse(SError.InValidEmailError, "无效邮箱"))
		c.Abort()
		return
	}

	name := c.PostForm("name")
	if !TypeValid.CheckNameValid(name) {
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
// @Success 200 {object} util.Response
// @Router /user/delete [delete]
func DeleteUser(c *gin.Context) {
	user := model.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	if !TypeValid.CheckIdExist(id) {
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
// @Success 200 {object} util.Response
// @Router /user/update [put]
func UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	if !TypeValid.CheckIdExist(id) {
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
	if TypeValid.CheckNameValid(name) {
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

var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// websocket test:https://www.qianbo.com.cn/Tool/WebSocket/

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("websocket upgrade failed: %v", err)
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Printf("close websocket conn failed: %v", err)
		}
	}(ws)

	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	msg, err := model.Subscribe(c, model.PublishKey)
	if err != nil {
		fmt.Printf("Redis Subscribe failed: %v", err)
	}

	formatTime := time.Now().Format("2006-01-02 15:04:05")
	m := fmt.Sprintf("[ws][%s]:%s", formatTime, msg)
	err = ws.WriteMessage(1, []byte(m))
	if err != nil {
		fmt.Printf("write message failed: %v", err)
	}
}
