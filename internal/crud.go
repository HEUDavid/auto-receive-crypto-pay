package internal

import (
	"context"
	"gorm.io/gorm"
)

func GetTokenByAddress(c context.Context, db *gorm.DB, fromAddress string) ([]*ReceiptData, error) {
	var dataList []*ReceiptData
	if err := db.Table((&ReceiptData{}).TableName()).
		Omit("raw_data").
		Where("from_address = ?", fromAddress).
		Find(&dataList).Error; err != nil {
		return nil, err
	}
	return dataList, nil
}

func GetTokenDetails(c context.Context, db *gorm.DB, token string) (*ReceiptData, error) {
	data := &ReceiptData{}
	if err := db.Table((&ReceiptData{}).TableName()).
		Omit("raw_data").
		Where("token = ?", token).
		First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
