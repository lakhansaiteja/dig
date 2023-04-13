package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HistoryRecord struct {
	Request   DNSRequest
	Timestamp string
}

func history(c *gin.Context) {
	var history []HistoryRecord
	mu.Lock() // locking to prevent race condition
	for e, count := dll.Back(), 0; e != nil && count < 30; e, count = e.Prev(), count+1 {
		req, ok := e.Value.(HistoryRecord)
		if ok {
			history = append(history, req)
		}
	}
	mu.Unlock()
	c.JSON(http.StatusOK, history)
}
