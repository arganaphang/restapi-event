package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"log"

	"github.com/Shopify/sarama"
)

const (
	TOPIC = "user_created"
)

func ConnectConsumer() (sarama.ConsumerGroup, error) {
	urls := os.Getenv("BROKER_URLS")
	if urls == "" {
		urls = "localhost:19092"
	}
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	conn, err := sarama.NewConsumerGroup(strings.Split(urls, ","), "user_created_consumer", config)
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
	keepRunning := true
	ctx, cancel := context.WithCancel(context.Background())
	worker, err := ConnectConsumer()
	if err != nil {
		log.Fatal("Failed to connect into stream")
	}
	consumer := Consumer{
		ready: make(chan bool),
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
	cancel()
	wg.Wait()
	if err = worker.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
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
	for {
		select {
		case message := <-claim.Messages():
			// TODO: Insert into database
			fmt.Println(string(message.Value))
			session.MarkMessage(message, "created")
		case <-session.Context().Done():
			return nil
		}
	}
}
