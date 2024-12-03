package model

import "gorm.io/gorm"

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
