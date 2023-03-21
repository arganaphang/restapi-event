package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
)

const (
	PORT  = 8000
	TOPIC = "user_created"
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

func ConnectProducer() (sarama.SyncProducer, error) {
	urls := os.Getenv("BROKER_URLS")
	if urls == "" {
		urls = "localhost:19092"
	}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	return sarama.NewSyncProducer(strings.Split(urls, ","), config)
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
		go func(data []Transaction) {
			producer, err := ConnectProducer()
			if err != nil {
				log.Println("Failed to connect into stream")
			}
			defer producer.Close()
			for _, trx := range data {
				messageByte, _ := json.Marshal(trx)
				msg := &sarama.ProducerMessage{
					Topic: TOPIC,
					Value: sarama.StringEncoder(string(messageByte)),
				}
				_, _, err := producer.SendMessage(msg)
				if err != nil {
					log.Println("Failed to push message")
				}
			}
		}(body.Data)
		ctx.JSON(http.StatusCreated, map[string]string{
			"message": "Transaction created",
		})
	})
	app.Run(fmt.Sprintf("0.0.0.0:%d", PORT))
}
