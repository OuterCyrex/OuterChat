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

func CoopCreate(tx *gorm.DB, OriginRequest Contact, OppositeRequest Contact) error {
	err := tx.Save(&OriginRequest).Error
	if err != nil {
		return err
	}
	err = tx.Create(&OppositeRequest).Error
	if err != nil {
		return err
	}
	return nil
}

func CoopDelete(tx *gorm.DB, OriginRequest Contact, OppositeRequest Contact) error {
	err := tx.Delete(&OriginRequest).Error
	if err != nil {
		return err
	}
	err = tx.Delete(&OppositeRequest).Error
	if err != nil {
		return err
	}
	return nil
}

func TransactionOperation(tx *gorm.DB, OriginRequest Contact, OppositeRequest Contact, opt func(*gorm.DB, Contact, Contact) error) error {
	transaction := tx.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	err := opt(transaction, OriginRequest, OppositeRequest)
	if err != nil {
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return nil
}

// Options pattern

type Options func(*gorm.DB) *gorm.DB

func WithStatus(status int) Options {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}

func FindContactByCoopId(db *gorm.DB, FromId uint, TargetId uint, opt ...Options) Contact {
	if FromId == TargetId {
		return Contact{}
	}
	for _, opt := range opt {
		db = opt(db)
	}
	var Result Contact
	db.Model(&Contact{}).Where("owner_id = ? and target_id = ?", FromId, TargetId).First(&Result)
	return Result
}

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

func PushFriendRequest(FromId uint, TargetId uint, Desc string) (Contact, error) {
	FindContactByCoopId(DB, FromId, TargetId)
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
	err := DB.Where("id = ?", ContactId).Find(&FriendRequest).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Contact{}, xerrors.Errorf("Invalid FriendRequest ID: %d", ContactId)
	} else if err != nil {
		return Contact{}, err
	}
	FriendRequest.Status = Status

	OppositeRequest := Contact{
		OwnerId:  FriendRequest.TargetId,
		TargetId: FriendRequest.OwnerId,
		Status:   Accept,
	}

	err = TransactionOperation(DB, FriendRequest, OppositeRequest, CoopCreate)
	if err != nil {
		return Contact{}, err
	}

	return FriendRequest, nil
}

func DeleteFriend(FromId uint, TargetId uint) error {
	contact, oppositeContact := Contact{}, Contact{}
	err := DB.Model(&Contact{}).Where("owner_id = ? and target_id = ? and status = ?", FromId, TargetId, Accept).First(&contact).Error
	if err != nil {
		return err
	}
	err = DB.Model(&Contact{}).Where("owner_id = ? and target_id = ? and status = ?", TargetId, FromId, Accept).First(&oppositeContact).Error
	if err != nil {
		return err
	}
	err = TransactionOperation(DB, contact, oppositeContact, CoopDelete)
	if err != nil {
		return err
	}
	return nil
}
