package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables for local testing
	err := godotenv.Load(".env")

	// Telegram Bot
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe("0.0.0.0:80", nil)

	// producer config
	config := &kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVER"), // kafka broker addr
	}

	// target topic name
	topic := "incoming_messages"

	// Create Kafka producer
	producer, err := kafka.NewProducer(config)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		log.Printf("%+v\n", update)
		value, err1 := json.Marshal(update)
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		log.Printf("Send message: %s", value)

		err2 := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(value),
		}, nil)

		if err2 != nil {
			fmt.Printf("Failed to produce message\n")
		} else {
			fmt.Printf("Produced message\n")
			producer.Flush(10)
		}
	}

	// Close Kafka producer
	producer.Close()
}
