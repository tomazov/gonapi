package main

import (
	"fmt"
	"log"
	"napi/internal/mq"
)

func main() {
	fmt.Println("🚀 Запуск RabbitMQ worker...")
	err := mq.StartConsumer()
	if err != nil {
		log.Fatalf("❌ Consumer error: %v", err)
	}
}
