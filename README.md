# Telegram echo bot using Kafka

With this repository I created a simple Telegram echo bot in order to learn more about Go and Kafka.

The echo bot consists of 3 services:
- telegram-callback: Endpoint for Telegram webhook. Writes content to incoming_messages Kafka topic.
- telegram-message-processor: Reads incoming_messages topic, processes the message and sent the client a message via outgoing_messages topic
- telegram-sender: Receives the outgoing_messages topic and creates the Telegram message

## Futher notes
- Run a service directly via `go run telegram-callback/telgram-callback.go`
- Set Telegram webhook via POST https://api.telegram.org/bot{{token}}/setWebhook?url=https://{{serveraddress}}/{{token}}
- The current setup is created for use in github actions only
- Kafka port docker internally: 29092
- Kafka port from outsite docker: 9092
- Keep your token safe!