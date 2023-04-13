package main

import (
	"container/list"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"sync"
)

var dll *list.List
var mu sync.Mutex

func init() {
	dll = list.New()
}

func main() {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.POST("/dig", dnsRecords)
		v1.GET("/dig/history", history)
	}

	// starting server on port 8080
	err := router.Run(":8080")
	if err != nil {
		log.Error(err, "server closed!")
	}
}
