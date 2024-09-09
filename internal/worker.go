package internal

import (
	"github.com/HEUDavid/auto-receive-crypto-pay/model"
	"github.com/HEUDavid/go-fsm/pkg"
	db "github.com/HEUDavid/go-fsm/pkg/db/mysql"
	mq "github.com/HEUDavid/go-fsm/pkg/mq/rmq"
	"github.com/HEUDavid/go-fsm/pkg/util"
	"sync"
)

type ServiceWorker struct {
	pkg.Worker[*ReceiptData]
}

func NewWorker() *ServiceWorker {
	w := &ServiceWorker{}
	w.MaxGoroutines = 10
	return w
}

var Worker = NewWorker()
var _initWorker sync.Once

func InitWorker() {
	_initWorker.Do(func() {
		Worker.RegisterModel(
			&ReceiptData{},
			&model.Task{},
			&model.UniqueRequest{},
		)
		Worker.RegisterFSM(ReceiptFSM)
		Worker.RegisterGenerator(util.UniqueID)
		Worker.RegisterDB(&db.Factory{Section: "mysql_public"})
		Worker.RegisterMQ(&mq.Factory{Section: "rmq_public"})
		Worker.Config = util.GetConfig()
		Worker.Init()
	})
}
