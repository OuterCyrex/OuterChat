package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name          string    `json:"name" gorm:"unique"`
	Password      string    `json:"password"`
	Phone         string    `json:"phone" gorm:"type:varchar(200)"`
	Email         string    `json:"email" valid:"email" gorm:"type:varchar(200)"`
	ClientIP      string    `json:"client_ip" gorm:"type:varchar(100)"`
	ClientPort    string    `json:"client_port" gorm:"type:varchar(50)"`
	LoginTime     time.Time `json:"login_time"`
	HeartbeatTime time.Time `json:"heartbeat_time"`
	LogoutTime    time.Time `json:"logout_time"`
	IsLogout      bool      `json:"is_logout"`
	DeviceInfo    string    `json:"device_info"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

// 用户列表查询

func GetUserList() []UserBasic {
	var userList []UserBasic
	err := DB.Model(&UserBasic{}).Find(&userList).Error
	if err != nil {
		fmt.Println("用户查询出错:", err)
	}
	return userList
}

func FindUserByField(Field string, data interface{}) (UserBasic, error) {
	user := UserBasic{}
	typeOf := reflect.TypeOf(&user).Elem()

	if _, ok := typeOf.FieldByName(Field); ok {
		err := DB.Where(fmt.Sprintf("%s = ?", Field), data).First(&user).Error
		if err != nil {
			return UserBasic{}, err
		}
		return user, nil
	}
	return UserBasic{}, errors.New("invalid field:" + Field)
}

func CreateUser(user UserBasic) *gorm.DB {
	return DB.Create(&user)
}
func DeleteUser(user UserBasic) *gorm.DB {
	return DB.Delete(&user)
}
func UpdateUser(user UserBasic) *gorm.DB {
	return DB.Model(&user).Where("id = ?", user.ID).Updates(UserBasic{Name: user.Name, Password: user.Password})
}

func LoginByName(username string, password string) *gorm.DB {
	return DB.Model(&UserBasic{}).Where("Name = ? and Password = ?", username, password).First(&UserBasic{})
}
