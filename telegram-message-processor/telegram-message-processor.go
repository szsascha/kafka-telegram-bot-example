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

	// Set up Kafka consumer configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVER"),                    // broker addr
		"group.id":          os.Getenv("KAFKA_GROUP_ID_TELEGRAM_MESSAGE_PROCESSOR"), // consumer group
		"auto.offset.reset": "earliest",                                             // earliest offset

	}

	// Create Kafka consumer
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		panic(err)
	}

	// subscribe to target topic
	consumer.SubscribeTopics([]string{"incoming_messages"}, nil)

	// Consume messages
	for {

		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Received message: %s\n", string(msg.Value))
		} else {
			fmt.Printf("Error while consuming message: %v (%v)\n", err, msg)
		}

		var tupdate tgbotapi.Update
		err1 := json.Unmarshal(json.RawMessage(msg.Value), &tupdate)
		if err1 != nil {
			fmt.Printf("Error reading message: %v\n", err1)
			continue
		}
		fmt.Printf("\n\n%+v\n\n", tupdate)
		outgoingMessage := OutgoingMessage{
			ChatId:  tupdate.Message.Chat.ID,
			Message: tupdate.Message.Text,
		}
		marshalledOutgoingMessage, err2 := json.Marshal(outgoingMessage)
		if err2 != nil {
			panic(err2)
		}

		// target topic name
		topic := "outgoing_messages"

		// Create Kafka producer
		producer, err3 := kafka.NewProducer(config)
		if err3 != nil {
			panic(err3)
		}

		log.Printf("Send message: %s", marshalledOutgoingMessage)

		err4 := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(marshalledOutgoingMessage),
		}, nil)

		if err4 != nil {
			fmt.Printf("Failed to produce message\n")
		} else {
			fmt.Printf("Produced message\n")
			producer.Flush(10)
		}

		// Close Kafka producer
		producer.Close()

	}

	// Close Kafka consumer
	consumer.Close()
}
