package model

import (
	"errors"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	OwnerId  uint
	TargetId uint
	Status   int
	Desc     string
}

const (
	Waiting = iota
	Accept
	Refused
)

func (table *Contact) TableName() string {
	return "contact"
}

// Options pattern

type Options func(*gorm.DB) *gorm.DB

func WithStatus(status int) Options {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}

func FindContactByCoopId(db *gorm.DB, FromId uint, TargetId uint, opt ...Options) (Contact, error) {
	if FromId == TargetId {
		return Contact{}, xerrors.Errorf("Same Ids for From and Target")
	}
	for _, opt := range opt {
		db = opt(db)
	}
	var Result Contact
	err := db.Model(&Contact{}).Where("owner_id = ? and target_id = ?", FromId, TargetId).First(&Result).Error
	return Result, err
}

// Others

func GetFriendListById(id uint) ([]UserBasic, error) {
	var contacts []Contact
	err := DB.Where("owner_id = ? and status = ?", id, Accept).Find(&contacts).Error
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

func IsFriendStatus(FromId uint, TargetId uint, opt Options) bool {
	if _, err := FindContactByCoopId(DB, FromId, TargetId, opt); errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil {
		return false
	}
	return true
}

func PushFriendRequest(FromId uint, TargetId uint, Desc string) (Contact, error) {
	WaitingRequest := Contact{
		OwnerId:  FromId,
		TargetId: TargetId,
		Status:   Waiting,
	}
	if Desc != "" {
		WaitingRequest.Desc = Desc
	}
	if err := DB.Create(&WaitingRequest).Error; err != nil {
		return Contact{}, err
	}
	return WaitingRequest, nil
}

// GetRequest has two Status, option == 1 means that to get Received Requests
// while option == 2 means that to get sent Requests
func GetRequest(id uint, option int) ([]Contact, error) {
	if !CheckIdExist(int(id)) {
		return nil, xerrors.Errorf("Invalid ID")
	}
	var queryString *gorm.DB
	switch option {
	case 1:
		queryString = DB.Model(&Contact{}).Where("status = ? and target_id = ?", Waiting, id)
	case 2:
		queryString = DB.Model(&Contact{}).Where("status = ? and owner_id = ?", Waiting, id)
	default:
		return nil, xerrors.Errorf("Invalid option")
	}
	var Result []Contact
	err := queryString.Find(&Result).Error
	if err != nil {
		return nil, err
	}
	return Result, nil
}

func DealWithFriendRequest(ContactId uint, Status int) (Contact, error) {
	var FriendRequest Contact

	switch Status {
	case Accept:
	case Refused:
	default:
		return Contact{}, xerrors.Errorf("Invalid status: %d", Status)
	}

	err := DB.Model(&Contact{}).Where("id = ? and status = ?", ContactId, Waiting).First(&FriendRequest).Error
	if err != nil {
		return Contact{}, err
	}
	FriendRequest.Status = Status

	OppositeRequest := Contact{
		OwnerId:  FriendRequest.TargetId,
		TargetId: FriendRequest.OwnerId,
		Status:   Accept,
	}

	tx := DB.Begin()
	err = tx.Save(&FriendRequest).Error
	if err != nil {
		tx.Rollback()
		return Contact{}, err
	}
	err = tx.Create(&OppositeRequest).Error
	if err != nil {
		tx.Rollback()
		return Contact{}, err
	}
	err = tx.Commit().Error
	if err != nil {
		return Contact{}, err
	}

	return FriendRequest, nil
}

func DeleteFriend(FromId uint, TargetId uint) error {
	contact, err := FindContactByCoopId(DB, FromId, TargetId, WithStatus(Accept))
	if err != nil {
		return err
	}
	oppositeContact, err := FindContactByCoopId(DB, TargetId, FromId, WithStatus(Accept))
	if err != nil {
		return err
	}

	tx := DB.Begin()
	err = tx.Delete(&contact).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Delete(&oppositeContact).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}
