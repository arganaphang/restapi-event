package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	PORT = 8000
)

type Transaction struct {
	ID        uint64    `json:"id"`
	Customer  string    `json:"customer"`
	Quantity  uint16    `json:"quantity"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

type RequestBody struct {
	RequestID uint64        `json:"request_id"`
	Data      []Transaction `json:"data" binding:"dive"`
}

func main() {
	app := gin.New()
	app.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "OK",
		})
	})
	app.POST("/save", func(ctx *gin.Context) {
		var body RequestBody
		if err := ctx.Bind(&body); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, map[string]string{
				"message": "Failed to created Transaction",
			})
			return
		}
		ctx.JSON(http.StatusCreated, map[string]string{
			"message": "Transaction created",
		})
	})
	app.Run(fmt.Sprintf("0.0.0.0:%d", PORT))
}
