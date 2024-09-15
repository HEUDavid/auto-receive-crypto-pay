package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/HEUDavid/auto-receive-crypto-pay/internal/parser"
	"github.com/HEUDavid/auto-receive-crypto-pay/model"
	. "github.com/HEUDavid/go-fsm/pkg/metadata"
	"log"
	"time"
)

var (
	// Hook 监听地址活动，收到 Node Services 回调, 节点服务参考：https://ethereum.org/en/developers/docs/nodes-and-clients/nodes-as-a-service/#popular-node-services
	// Processed 回调请求落库之后，执行自定义逻辑
	Hook      = GenState("Hook", false, hookHandler)
	Processed = State[*ReceiptData]{Name: "Processed", IsFinal: true, Handler: nil}

	// New 为"发送地址FromAddress"，启动执行流程，可以自行拓展
	New        = GenState("New", false, newHandler)
	GenInvoice = GenState("GenInvoice", false, genInvoiceHandler)
	End        = State[*ReceiptData]{Name: "End", IsFinal: true, Handler: nil}
)

var (
	Hook2Processed = GenTransition(Hook, Processed)
	New2GenInvoice = GenTransition(New, GenInvoice)
	GenInvoice2End = GenTransition(GenInvoice, End)
)

var ReceiptFSM = func() FSM[*ReceiptData] {
	fsm := GenFSM[*ReceiptData]("Receipt")

	fsm.RegisterState(Hook, Processed)
	fsm.RegisterTransition(Hook2Processed)

	fsm.RegisterState(New, GenInvoice, End)
	fsm.RegisterTransition(New2GenInvoice, GenInvoice2End)
	return fsm
}()

func hookHandler(task *Task[*ReceiptData]) error {
	log.Printf("[FSM] State: %s, Task.Data: %s", task.State, _pretty(task.GetData()))
	task.Data.Comment = "webhook payload"

	// log.Println(_pretty((*task.GetData()).RawData))
	var rawData parser.WebhookData
	if err := json.Unmarshal((*task.GetData()).RawData, &rawData); err != nil {
		fmt.Println("json.Unmarshal error: ", err)
	}

	// 检查数据，根据链上交易信息创建任务
	addressConfig, exists := GetConfig().AdminAddress[rawData.Event.Network]
	if !exists {
		task.Data.Comment = "no relevant network is configured"
		task.State = Processed.GetName()
		return nil
	}

	var adminAddress []string
	for _, a := range addressConfig {
		adminAddress = append(adminAddress, a.Address)
	}
	for _, a := range rawData.Event.Activity {
		if !_contains(adminAddress, a.ToAddress) {
			continue
		}

		task.Data.Comment = "process receipt"
		logicTask := GenTaskInstance(a.Hash, "", &ReceiptData{Data: model.Data{
			Network:         rawData.Event.Network,
			Hash:            a.Hash,
			FromAddress:     a.FromAddress,
			ToAddress:       a.ToAddress,
			Asset:           a.Asset,
			Value:           a.Value,
			RawData:         (*task.GetData()).RawData,
			TransactionTime: task.Data.TransactionTime,
		}})
		logicTask.Type = "Logic"
		logicTask.State = New.GetName()
		if err := Adapter.Create(context.Background(), logicTask); err != nil {
			log.Printf("[FSM] Create logic task error: %s", err)
			return err
		}

	}

	task.State = Processed.GetName() // 标记为已处理
	return nil
}

func newHandler(task *Task[*ReceiptData]) error {
	log.Printf("[FSM] State: %s, Task.Data: %s", task.State, _pretty(task.GetData()))

	// It may be necessary to perform some checks.
	// It may be necessary to pre-record the request to the database to ensure idempotency.
	// For example, generating some request IDs.
	// ...

	task.Data.Comment = "receive cryptocurrency"
	task.State = GenInvoice.GetName() // 下一步：譬如生成发票
	return nil
}

func genInvoiceHandler(task *Task[*ReceiptData]) error {
	log.Printf("[FSM] State: %s, Task: %s", task.State, _pretty(task))

	// Invoke RPC interfaces to perform certain operations.
	// Here we generate the ID directly.
	task.Data.InvoiceID = Worker.GenID()

	currentTime := time.Now()
	timeAfter30Days := currentTime.AddDate(0, 0, 30)
	task.Data.ValidFrom = uint64(currentTime.Unix())
	task.Data.ValidTo = uint64(timeAfter30Days.Unix())

	// Maybe some notifications need to be sent
	log.Println("send mail...")

	task.Data.Comment = "gen invoice success"
	task.State = End.GetName() // Switch to next state
	return nil
}

func _pretty(v interface{}) string {
	s, _ := json.MarshalIndent(v, "", "  ")
	return string(s)
}

func _contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}
