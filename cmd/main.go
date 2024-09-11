package main

import (
	. "github.com/HEUDavid/auto-receive-crypto-pay/internal"
	"github.com/HEUDavid/auto-receive-crypto-pay/model"
	. "github.com/HEUDavid/go-fsm/pkg/metadata"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

func Webhook(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	body, err := io.ReadAll(c.Request.Body)
	log.Printf("webhook payload: %s, %v\n", body, err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "no data"})
		return
	}

	task := GenTaskInstance(Adapter.GenID(), "", &ReceiptData{Data: model.Data{RawData: body}})
	task.Type = "Persist"
	task.State = Hook.GetName()

	_response(c, Adapter.Create(c, task), task)
}

func Query(c *gin.Context) {
	task := GenTaskInstance(c.Query("request_id"), c.Query("task_id"), &ReceiptData{})
	_response(c, Adapter.Query(c, task), task)
}

func QueryToken(c *gin.Context) {
	data, err := GetToken(c, Adapter.GetDB(), c.Query("from_address"))
	_response(c, err, struct {
		Token     string
		ValidFrom uint64
		ValidTo   uint64
	}{
		data.Token,
		data.ValidFrom,
		data.ValidTo,
	})
}

func _response(c *gin.Context, err error, task interface{}) {
	if err == nil {
		c.JSON(http.StatusOK, task)
	} else {
		c.JSON(http.StatusOK, map[string]string{"error": err.Error()})
	}
}

func init() {
	gin.SetMode(GetConfig().Global.Mode)

	f, _ := os.Create(GetConfig().Global.LogPath)
	mw := io.MultiWriter(os.Stdout, f)
	gin.DefaultWriter = mw
	gin.DefaultErrorWriter = mw
	log.SetOutput(gin.DefaultWriter)

	InitWorker()
	InitAdapter()
}

func main() {
	r := gin.Default()
	r.POST("/webhook", Webhook)
	r.GET("/query", Query)
	r.GET("/query_token", QueryToken)

	Worker.Run()
	log.Println("[FSM] Worker started...")

	_ = r.Run(GetConfig().Global.Addr)
}
