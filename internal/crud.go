package internal

import (
	"context"
	"gorm.io/gorm"
)

func GetToken(c context.Context, db *gorm.DB, fromAddress string) ([]*ReceiptData, error) {
	var dataList []*ReceiptData
	if err := db.Table((&ReceiptData{}).TableName()).
		Select("*").Omit("raw_data").
		Where("from_address = ?", fromAddress).
		Find(&dataList).Error; err != nil {
		return nil, err
	}
	return dataList, nil
}
