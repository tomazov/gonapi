package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"napi/pkg/config"
	"napi/internal/api"
	"napi/internal/mq"
	"napi/internal/cache"
)

func main() {
	// 1. Load env config
	config.Load()

	// 2. Init RabbitMQ connection
	if err := mq.Init(config.Cfg.RabbitMQ); err != nil {
		log.Fatalf("🐇 Failed to init RabbitMQ: %v", err)
	}
	defer mq.Close()

	// 3. Init Memcached client
	if err := cache.Init(config.Cfg.MemcachedURL); err != nil {
		log.Fatalf("🧊 Failed to connect to Memcached: %v", err)
	}

	// 4. Start Fiber app with config
	app := fiber.New(fiber.Config{
		ReadTimeout: config.Cfg.ServerReadTimeout,
	})

	// 5. Register API routes
	api.SetupRoutes(app)

	// 6. Run server
	log.Printf("🚀 Server started on %s (env: %s)", config.Cfg.ServerURL, config.Cfg.AppEnv)
	if err := app.Listen(config.Cfg.ServerURL); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
