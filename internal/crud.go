package internal

import (
	. "context"
	"gorm.io/gorm"
)

func QueryToken(c Context, db *gorm.DB, fromAddress string) (*ReceiptData, error) {
	data := &ReceiptData{}
	if err := db.Table(data.TableName()).Where("from_address = ?", fromAddress).Find(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
