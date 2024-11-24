package internal

import (
	"fmt"
	"github.com/HEUDavid/auto-receive-crypto-pay/model"
	. "github.com/HEUDavid/go-fsm/pkg/metadata"
	"github.com/HEUDavid/go-fsm/pkg/util"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
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

	response(c, Adapter.Create(c, task), task)
}

func QueryTask(c *gin.Context) {
	task := GenTaskInstance(c.Query("request_id"), c.Query("task_id"), &ReceiptData{})
	response(c, Adapter.Query(c, task), task)
}

type invoice struct {
	Network         string
	FromAddress     string
	ToAddress       string
	Asset           string
	Value           float64
	InvoiceID       string
	ValidFrom       time.Time
	ValidTo         time.Time
	TransactionTime time.Time
	Content         string
	Valid           bool
}

func isValid(validFrom, validTo uint64) bool {
	currentTime := uint64(time.Now().Unix())
	return currentTime >= validFrom && currentTime <= validTo
}

func toInvoice(data *ReceiptData) *invoice {
	if data == nil {
		return nil
	}
	return &invoice{
		Network:         data.Network,
		FromAddress:     data.FromAddress,
		ToAddress:       data.ToAddress,
		Asset:           data.Asset,
		Value:           data.Value,
		InvoiceID:       data.InvoiceID,
		ValidFrom:       time.Unix(int64(data.ValidFrom), 0),
		ValidTo:         time.Unix(int64(data.ValidTo), 0),
		TransactionTime: time.Unix(int64(data.TransactionTime), 0),
		Content:         data.InvoiceID,
		Valid:           isValid(data.ValidFrom, data.ValidTo),
	}
}

func checkRequest(c *gin.Context, required ...string) (bool, string) {
	for _, param := range required {
		if c.Query(param) == "" {
			return false, param
		}
	}
	return true, ""
}

func QueryInvoice(c *gin.Context) {
	if valid, missingParam := checkRequest(c, "from_address"); !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing required parameter: " + missingParam,
		})
		return
	}

	dataList, err := GetInvoiceByAddress(c, Adapter.GetDB(), c.Query("from_address"))
	var invoices []*invoice
	for _, data := range dataList {
		invoices = append(invoices, toInvoice(data))
	}
	response(c, err, invoices)
}

func InvoiceDetails(c *gin.Context) {
	if valid, missingParam := checkRequest(c, "invoice"); !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing required parameter: " + missingParam,
		})
		return
	}

	data, err := GetInvoiceDetails(c, Adapter.GetDB(), c.Query("invoice"))
	response(c, err, toInvoice(data))
}

func response(c *gin.Context, err error, task interface{}) {
	if err == nil {
		c.JSON(http.StatusOK, task)
	} else {
		c.JSON(http.StatusOK, map[string]string{"error": err.Error()})
	}
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"hostRoot":       GetConfig().Global.HostRoot,
		"adminAddresses": GetConfig().AdminAddress,
	})
}

func Router(path string) string {
	return fmt.Sprintf("%s/%s", GetConfig().Global.HostRoot, path)
}

func Source(path string) string {
	return fmt.Sprintf("%s/%s", util.FindProjectRoot(), path)
}
