package mq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"napi/internal/workers"
	"os"
)

var ConsumerChannel *amqp.Channel

type MQPayload struct {
	SearchID string `json:"search_id"`
	RecID    int    `json:"rec_id"` // => operator ID
	Payload  any    `json:"payload"`
}

func InitConsumer(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel for consumer: %w", err)
	}

	ConsumerChannel = ch
	return nil
}

func StartConsumer(queueName string) error {
	if ConsumerChannel == nil {
		return fmt.Errorf("consumer channel not initialized")
	}

	msgs, err := ConsumerChannel.Consume(
		queueName,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume queue '%s': %w", queueName, err)
	}

	log.Printf("📡 Слухаємо чергу: %s", queueName)

	for msg := range msgs {
		var task MQPayload
		if err := json.Unmarshal(msg.Body, &task); err != nil {
			log.Printf("❌ Error decoding msg: %v", err)
			continue
		}

		log.Printf("📥 Отримано завдання: rec_id=%d searchId=%s", task.RecID, task.SearchID)

		// запуск обробника адаптера
		if err := workers.RunAdapter(task.RecID, task.SearchID, task.Payload); err != nil {
			log.Printf("❌ Adapter error: %v", err)
		}
	}

	return nil
}

func GetAMQPConn() (*amqp.Connection, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%s%s",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
		os.Getenv("RABBITMQ_VHOST"),
	)
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	return conn, nil
}
