package main

import (
	"fmt"
	. "github.com/HEUDavid/auto-receive-crypto-pay/internal"
	"github.com/HEUDavid/auto-receive-crypto-pay/model"
	. "github.com/HEUDavid/go-fsm/pkg/metadata"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func Webhook(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	body, err := io.ReadAll(c.Request.Body)
	fmt.Printf("webhook payload: %s, %v\n", body, err)
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

func _response(c *gin.Context, err error, task interface{}) {
	if err == nil {
		c.JSON(http.StatusOK, task)
	} else {
		c.JSON(http.StatusOK, map[string]string{"error": err.Error()})
	}
}

func init() {
	InitWorker()
	InitAdapter()
}

func main() {
	Worker.Run()
	log.Println("[FSM] Worker started...")

	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.POST("/webhook", Webhook)
	r.GET("/query", Query)
	_ = r.Run("127.0.0.1:8080")
}
