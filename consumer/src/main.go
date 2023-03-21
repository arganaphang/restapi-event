package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Shopify/sarama"
)

const (
	TOPIC    = "user_created"
	GROUP_ID = "user_created_consumer"
)

func ConnectConsumer() (sarama.ConsumerGroup, error) {
	urls := os.Getenv("BROKER_URLS")
	if urls == "" {
		urls = "localhost:19092"
	}
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	conn, err := sarama.NewConsumerGroup(strings.Split(urls, ","), GROUP_ID, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type Transaction struct {
	ID        uint64    `json:"id"`
	Customer  string    `json:"customer"`
	Quantity  uint16    `json:"quantity"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgresql://postgres:mystrongpassword@localhost:5432/restapi?sslmode=disable"
	}
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		log.Fatalln(err)
	}
	keepRunning := true
	ctx, cancel := context.WithCancel(context.Background())
	worker, err := ConnectConsumer()
	if err != nil {
		log.Fatal("Failed to connect into stream")
	}
	consumer := Consumer{
		ready: make(chan bool),
		db:    db,
	}
	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := worker.Consume(ctx, strings.Split(TOPIC, ","), &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")
	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(worker, &consumptionIsPaused)
		}
	}
	// ? Gracefully shutdown
	cancel()
	wg.Wait()
	if err = worker.Close(); err != nil {
		log.Panicf("Error closing Consumer: %v", err)
	}
	if err = db.Close(); err != nil {
		log.Panicf("Error closing database: %v", err)
	}
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
	db    *sqlx.DB
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	count := 1
	for {
		select {
		case msg := <-claim.Messages():
			var trx Transaction
			if err := json.Unmarshal(msg.Value, &trx); err != nil {
				log.Println("Failed to serialize data") // ? <- this should never gonna happened
				continue
			}
			_, err := consumer.db.Exec(`INSERT INTO "public"."transactions" ("id", "customer", "quantity", "price", "timestamp") VALUES ($1, $2, $3, $4, $5)`, trx.ID, trx.Customer, trx.Quantity, trx.Quantity, trx.Timestamp)
			if err != nil {
				log.Println("Failed to insert data ", err)
				continue
			}
			log.Println(count)
			count++
			session.MarkMessage(msg, "created")
		case <-session.Context().Done():
			return nil
		}
	}
}
