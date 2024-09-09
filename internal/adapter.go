package internal

import (
	"context"
	"github.com/HEUDavid/auto-receive-crypto-pay/model"
	"github.com/HEUDavid/go-fsm/pkg"
	db "github.com/HEUDavid/go-fsm/pkg/db/mysql"
	. "github.com/HEUDavid/go-fsm/pkg/metadata"
	mq "github.com/HEUDavid/go-fsm/pkg/mq/rmq"
	"github.com/HEUDavid/go-fsm/pkg/util"
	"log"
	"sync"
	"time"
)

type ServiceAdapter struct {
	pkg.Adapter[*ReceiptData]
}

func (a *ServiceAdapter) BeforeCreate(c context.Context, task *Task[*ReceiptData]) error {
	log.Println("[FSM] Rewrite BeforeCreate...")
	task.Version = 1
	task.Data.TransactionTime = uint64(time.Now().Unix())
	return nil
}

func NewAdapter() *ServiceAdapter {
	a := &ServiceAdapter{}
	a.ReBeforeCreate = a.BeforeCreate
	return a
}

var Adapter = NewAdapter()
var _initAdapter sync.Once

func InitAdapter() {
	_initAdapter.Do(func() {
		Adapter.RegisterModel(
			&ReceiptData{},
			&model.Task{},
			&model.UniqueRequest{},
		)
		Adapter.RegisterFSM(ReceiptFSM)
		Adapter.RegisterGenerator(util.UniqueID)
		Adapter.RegisterDB(&db.Factory{Section: "mysql_public"})
		Adapter.RegisterMQ(&mq.Factory{Section: "rmq_public"})
		Adapter.Config = util.GetConfig()
		_ = Adapter.Init()
	})
}
