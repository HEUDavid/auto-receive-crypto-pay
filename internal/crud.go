package internal

import (
	. "context"
	"gorm.io/gorm"
)

func GetToken(c Context, db *gorm.DB, fromAddress string) (*ReceiptData, error) {
	data := &ReceiptData{}
	if err := db.Table(data.TableName()).Select("token, valid_from, valid_to").
		Where("from_address = ?", fromAddress).Find(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
