package main

import (
	"fmt"
	. "github.com/HEUDavid/auto-receive-crypto-pay/internal"
	"github.com/HEUDavid/go-fsm/pkg/util"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"path/filepath"
)

func init() {
	gin.SetMode(GetConfig().Global.Mode)

	logPath := filepath.Join(util.FindProjectRoot(), GetConfig().Global.LogPath)
	if err := os.MkdirAll(filepath.Dir(logPath), os.ModePerm); err != nil {
		panic(fmt.Sprintf("Failed to create log directory: %v", err))
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open or create log file: %v", err))
	}

	mw := io.MultiWriter(os.Stdout, f)
	gin.DefaultWriter = mw
	log.SetOutput(mw)

	InitWorker()
	InitAdapter()
}

func main() {
	r := gin.Default()

	r.POST(Router("webhook"), Webhook)
	r.GET(Router("query"), QueryTask)
	r.GET(Router("query_token"), QueryInvoice)
	r.GET(Router("token_details"), InvoiceDetails)

	r.GET(Router("pay"), Index)

	r.Static(Router("src"), Source("static/src"))
	r.LoadHTMLGlob(Source("static/templates/*"))

	Worker.Run()
	log.Println("[FSM] Worker started...")

	log.Printf("[SERVICE] Listening on %s%s", GetConfig().Global.Addr, GetConfig().Global.HostRoot)
	_ = r.Run(GetConfig().Global.Addr)
}
