package mq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

var ProducerChannel *amqp.Channel

func InitProducer(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel for producer: %w", err)
	}

	ProducerChannel = ch
	return nil
}

func PublishToQueue(queueName string, data any) error {
	if ProducerChannel == nil {
		return fmt.Errorf("producer channel not initialized")
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = ProducerChannel.Publish(
		"",         // exchange
		queueName,  // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("📤 Published message to queue '%s': %s", queueName, body)
	return nil
}
