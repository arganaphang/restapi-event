package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	databaseURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DATABASE_HOST"),
		"5432",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
	)
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		log.Fatalln(err)
	}
	app := gin.New()
	app.SetTrustedProxies(nil)
	app.GET("/healthz", func(ctx *gin.Context) {
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
		// TODO: Insert To Database
		for _, trx := range body.Data {
			_, err := db.Exec(`INSERT INTO "public"."transactions" ("id", "customer", "quantity", "price", "timestamp") VALUES ($1, $2, $3, $4, $5)`, trx.ID, trx.Customer, trx.Quantity, trx.Price, trx.Timestamp)
			if err != nil {
				log.Println("Failed to insert data ", err)
			}
		}

		ctx.JSON(http.StatusCreated, map[string]string{
			"message": "Transaction created",
		})
	})
	app.Run(fmt.Sprintf("0.0.0.0:%d", PORT))
}
