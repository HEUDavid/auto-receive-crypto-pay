package internal

import (
	"context"
	"gorm.io/gorm"
)

func GetInvoiceByAddress(c context.Context, db *gorm.DB, fromAddress string) ([]*ReceiptData, error) {
	var dataList []*ReceiptData
	if err := db.Table((&ReceiptData{}).TableName()).
		Omit("raw_data").
		Where("from_address = ?", fromAddress).
		Find(&dataList).Error; err != nil {
		return nil, err
	}
	return dataList, nil
}

func GetInvoiceDetails(c context.Context, db *gorm.DB, invoiceID string) (*ReceiptData, error) {
	data := &ReceiptData{}
	if err := db.Table((&ReceiptData{}).TableName()).
		Omit("raw_data").
		Where("invoice_id = ?", invoiceID).
		First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
