package internal

import (
	"encoding/json"
	. "github.com/HEUDavid/go-fsm/pkg/metadata"
	"log"
)

var (
	New  = GenState("New", false, newHandler)
	Send = GenState("Send", false, sendHandler)
	End  = State[*ReceiptData]{Name: "End", IsFinal: true, Handler: nil}
)

var (
	New2Send = GenTransition(New, Send)
	Send2End = GenTransition(Send, End)
	End2End  = GenTransition(End, End)
)

var ReceiptFSM = func() FSM[*ReceiptData] {
	fsm := GenFSM("ReceiptFSM", New)
	fsm.RegisterState(New, Send, End)
	fsm.RegisterTransition(New2Send, Send2End, End2End)
	return fsm
}()

func newHandler(task *Task[*ReceiptData]) error {
	log.Printf("[FSM] State: %s, Task.Data: %s", task.State, _pretty(task.GetData()))
	task.Data.Comment = "receive cryptocurrency"
	task.State = Send.GetName() // 下一步 发送邮件通知
	return nil
}

func sendHandler(task *Task[*ReceiptData]) error {
	log.Printf("[FSM] State: %s, Task: %s", task.State, _pretty(task))

	// Invoke RPC interfaces to perform certain operations.
	log.Println("send mail success")
	task.Data.Comment = "send mail success"

	task.State = End.GetName() // Switch to next state
	return nil
}

func _pretty(v interface{}) string {
	s, _ := json.MarshalIndent(v, "", "  ")
	return string(s)
}
