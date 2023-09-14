package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type OutgoingMessage struct {
	ChatId  int64
	Message string
}

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

	// Set up Kafka consumer configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVER"),         // broker addr
		"group.id":          os.Getenv("KAFKA_GROUP_ID_TELEGRAM_SENDER"), // consumer group
		"auto.offset.reset": "earliest",                                  // earliest offset

	}

	// Create Kafka consumer
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		panic(err)
	}

	// subscribe to target topic
	consumer.SubscribeTopics([]string{"outgoing_messages"}, nil)

	// Consume messages
	for {
		msg, err1 := consumer.ReadMessage(-1)
		var outgoingMessage OutgoingMessage
		if err1 == nil {
			err2 := json.Unmarshal(msg.Value, &outgoingMessage)
			if err2 != nil {
				panic(err2)
			}
			fmt.Printf("Received message: %s\n", string(msg.Value))
			tmsg := tgbotapi.NewMessage(outgoingMessage.ChatId, outgoingMessage.Message)
			if _, err1 := bot.Send(tmsg); err1 != nil {
				panic(err1)
			}
		} else {
			fmt.Printf("Error while consuming message: %v (%v)\n", err, msg)
		}

	}

	// Close Kafka consumer
	consumer.Close()
}
