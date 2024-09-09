package internal

import (
	"encoding/json"
	. "github.com/HEUDavid/go-fsm/pkg/metadata"
	"log"
)

var (
	Hook      = GenState("Hook", false, hookHandler)
	Processed = State[*ReceiptData]{Name: "Processed", IsFinal: true, Handler: nil}

	New      = GenState("New", false, newHandler)
	GenToken = GenState("GenToken", false, genTokenHandler)
	End      = State[*ReceiptData]{Name: "End", IsFinal: true, Handler: nil}
)

var (
	Hook2Processed = GenTransition(Hook, Processed)
	New2GenToken   = GenTransition(New, GenToken)
	GenToken2End   = GenTransition(GenToken, End)
)

var ReceiptFSM = func() FSM[*ReceiptData] {
	fsm := GenFSM[*ReceiptData]("Receipt")

	fsm.RegisterState(Hook, Processed)
	fsm.RegisterTransition(Hook2Processed)

	fsm.RegisterState(New, GenToken, End)
	fsm.RegisterTransition(New2GenToken, GenToken2End)
	return fsm
}()

func hookHandler(task *Task[*ReceiptData]) error {
	log.Printf("[FSM] State: %s, Task.Data: %s", task.State, _pretty(task.GetData()))
	task.Data.Comment = "webhook payload"
	// 检查数据，然后根据合法收据建立新的任务
	// ...
	task.State = Processed.GetName()
	return nil
}

func newHandler(task *Task[*ReceiptData]) error {
	log.Printf("[FSM] State: %s, Task.Data: %s", task.State, _pretty(task.GetData()))
	task.Data.Comment = "receive cryptocurrency"
	task.State = GenToken.GetName() // 下一步 执行一些动作
	return nil
}

func genTokenHandler(task *Task[*ReceiptData]) error {
	log.Printf("[FSM] State: %s, Task: %s", task.State, _pretty(task))

	// Invoke RPC interfaces to perform certain operations.
	// 生成卡密或者发送商品之类的
	log.Println("send mail...")
	task.Data.Comment = "send mail success"

	task.State = End.GetName() // Switch to next state
	return nil
}

func _pretty(v interface{}) string {
	s, _ := json.MarshalIndent(v, "", "  ")
	return string(s)
}
