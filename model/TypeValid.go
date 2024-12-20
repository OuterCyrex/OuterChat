package model

import (
	"OuterChat/util"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

func CheckUserValidByField(Field string, data interface{}) (bool, error) {
	user := UserBasic{}
	typeOf := reflect.TypeOf(&user).Elem()

	if _, ok := typeOf.FieldByName(Field); ok {
		err := DB.Where(fmt.Sprintf("%s = ?", Field), data).First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		return false, err
	}
	return false, errors.New("invalid field:" + Field)
}

func CheckIdExist(id int) bool {
	var user UserBasic
	user.ID = uint(id)
	err := DB.Where("ID = ?", id).First(&user).Error
	if err != nil {
		return false
	}
	return true
}

func CheckNameValid(name string) bool {
	ok, err := CheckUserValidByField("Name", name)
	if err != nil {
		fmt.Println("err: ", err)
		return false
	}
	return ok
}

func CheckEmailValid(email string) bool {
	ok, err := CheckUserValidByField("Email", email)
	if err != nil {
		fmt.Println("err: ", err)
		return false
	}
	return ok
}

func CheckTokenValid(tokenString string) bool {
	claims, err := util.ParseToken(tokenString)
	if err != nil {
		fmt.Printf("Parse Token Failed: %v", err)
		return false
	}
	return CheckIdExist(int(claims.UID))
}
