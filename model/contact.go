package model

import (
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	OwnerId  uint
	TargetId uint
	Type     int
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

func GetFriendListById(id uint) ([]UserBasic, error) {
	var contacts []Contact
	err := DB.Where("owner_id = ? and type = 1", id).Find(&contacts).Error
	if err != nil {
		return []UserBasic{}, err
	}
	var FriendList []UserBasic
	var Ids []uint
	for _, v := range contacts {
		Ids = append(Ids, v.TargetId)
	}
	err = DB.Where("id in (?)", Ids).Find(&FriendList).Error
	if err != nil {
		return []UserBasic{}, err
	}
	return FriendList, nil
}
