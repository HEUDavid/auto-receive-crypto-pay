package main

import (
	"fmt"
	. "github.com/HEUDavid/auto-receive-crypto-pay/internal"
	"github.com/HEUDavid/auto-receive-crypto-pay/model"
	. "github.com/HEUDavid/go-fsm/pkg/metadata"
	"github.com/HEUDavid/go-fsm/pkg/util"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func Webhook(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	if c.Query("auth") != GetConfig().Global.Auth {
		c.JSON(http.StatusBadRequest, gin.H{"error": "auth failed"})
		return
	}

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

type token struct {
	Token           string
	Valid           bool
	ValidFrom       time.Time
	ValidTo         time.Time
	Network         string
	FromAddress     string
	ToAddress       string
	Asset           string
	Value           float64
	TransactionTime time.Time
}

func isValid(validFrom, validTo uint64) bool {
	currentTime := uint64(time.Now().Unix())
	return currentTime >= validFrom && currentTime <= validTo
}

func toToken(data *ReceiptData) *token {
	return &token{
		Token:           data.Token,
		Valid:           isValid(data.ValidFrom, data.ValidTo),
		ValidFrom:       time.Unix(int64(data.ValidFrom), 0),
		ValidTo:         time.Unix(int64(data.ValidTo), 0),
		Network:         data.Network,
		FromAddress:     data.FromAddress,
		ToAddress:       data.ToAddress,
		Asset:           data.Asset,
		Value:           data.Value,
		TransactionTime: time.Unix(int64(data.TransactionTime), 0),
	}
}

func QueryToken(c *gin.Context) {
	dataList, err := GetTokenByAddress(c, Adapter.GetDB(), c.Query("from_address"))
	var tokens []*token
	for _, data := range dataList {
		tokens = append(tokens, toToken(data))
	}
	_response(c, err, tokens)
}

func TokenDetails(c *gin.Context) {
	data, err := GetTokenDetails(c, Adapter.GetDB(), c.Query("token"))
	_response(c, err, toToken(data))
}

func _response(c *gin.Context, err error, task interface{}) {
	if err == nil {
		c.JSON(http.StatusOK, task)
	} else {
		c.JSON(http.StatusOK, map[string]string{"error": err.Error()})
	}
}

func Index(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"hostRoot":       GetConfig().Global.HostRoot,
		"adminAddresses": GetConfig().AdminAddress,
	})
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

func router(path string) string {
	return fmt.Sprintf("%s/%s", GetConfig().Global.HostRoot, path)
}

func source(path string) string {
	return fmt.Sprintf("%s/%s", util.FindProjectRoot(), path)
}

func main() {
	r := gin.Default()

	r.POST(router("webhook"), Webhook)
	r.GET(router("query"), Query)
	r.GET(router("query_token"), QueryToken)
	r.GET(router("token_details"), TokenDetails)

	r.GET(router("pay"), Index)

	r.Static(router("src"), source("static/src"))
	r.LoadHTMLGlob(source("static/templates/*"))

	Worker.Run()
	log.Println("[FSM] Worker started...")

	log.Printf("Listening on %s%s", GetConfig().Global.Addr, GetConfig().Global.HostRoot)
	_ = r.Run(GetConfig().Global.Addr)
}
