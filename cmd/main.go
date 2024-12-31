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

func setupLog() {
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

	log.SetPrefix("[Auto Receive Crypto Pay] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func init() {
	gin.SetMode(GetConfig().Global.Mode)
	setupLog()
}

func main() {
	r := gin.Default()

	r.POST(Router("webhook"), Webhook)
	r.GET(Router("query_task"), QueryTask)

	r.GET(Router("pay"), Index)
	r.GET(Router("query_invoice"), QueryInvoice)
	r.GET(Router("invoice_details"), InvoiceDetails)

	r.Static(Router("src"), Source("static/src"))
	r.LoadHTMLGlob(Source("static/templates/*"))

	Worker.DoInit()
	Adapter.DoInit()

	Worker.Run()
	log.Println("[FSM] worker started...")

	log.Printf("[SERVICE] listening on %s%s", GetConfig().Global.Addr, GetConfig().Global.HostRoot)
	_ = r.Run(GetConfig().Global.Addr)

}
