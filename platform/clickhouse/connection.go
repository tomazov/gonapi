package clickhouse

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
)

var DB *sql.DB

func Connect() error {
	host := os.Getenv("CLICKHOUSE_HOST")
	db := os.Getenv("CLICKHOUSE_DB")
	user := os.Getenv("CLICKHOUSE_USER")
	pass := os.Getenv("CLICKHOUSE_PASSWORD")

	connStr := fmt.Sprintf("tcp://%s:9000?database=%s&username=%s&password=%s&read_timeout=10s&write_timeout=20s",
		host, db, user, pass,
	)

	var err error
	DB, err = sql.Open("clickhouse", connStr)
	if err != nil {
		return fmt.Errorf("failed to open ClickHouse connection: %w", err)
	}

	// Перевіримо, чи реально підʼєдналися
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	// Додаткові налаштування
	DB.SetMaxOpenConns(20)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(time.Minute * 5)

	log.Println("✅ ClickHouse connected successfully")
	return nil
}
