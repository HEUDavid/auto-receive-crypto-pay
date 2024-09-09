package internal

import "github.com/HEUDavid/auto-receive-crypto-pay/model"

type ReceiptData struct {
	model.Data
}

func (d *ReceiptData) SetTaskID(taskID string) {
	d.TaskID = taskID
}
