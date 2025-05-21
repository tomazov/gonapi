package rabbitmq

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

var Conn *amqp.Connection

func Connect() (*amqp.Connection, error) {
	if Conn != nil && !Conn.IsClosed() {
		return Conn, nil
	}

	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s%s",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
		os.Getenv("RABBITMQ_VHOST"),
	)

	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	Conn = conn
	log.Println("✅ RabbitMQ connection established")

	return conn, nil
}
