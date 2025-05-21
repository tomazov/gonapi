package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

var client *memcache.Client

const (
	defaultTTL = 900 // 15 хвилин
)

// Init підключає Memcached
func Init(address string) error {
	client = memcache.New(address)
	return client.Ping()
}

// SetJSON зберігає структуру у Memcached
func SetJSON(key string, value interface{}, ttl int32) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json.Marshal error: %w", err)
	}

	return client.Set(&memcache.Item{
		Key:        key,
		Value:      bytes,
		Expiration: ttl,
	})
}

// GetJSON отримує структуру з кешу
func GetJSON(key string, target interface{}) error {
	item, err := client.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(item.Value, target)
}

// SetString — простий string з TTL
func SetString(key string, val string, ttl int32) error {
	return client.Set(&memcache.Item{
		Key:        key,
		Value:      []byte(val),
		Expiration: ttl,
	})
}

// GetString — простий string
func GetString(key string) (string, error) {
	item, err := client.Get(key)
	if err != nil {
		return "", err
	}
	return string(item.Value), nil
}

// Delete ключа
func Delete(key string) error {
	return client.Delete(key)
}

// TTL за замовчуванням
func TTL() int32 {
	return int32(defaultTTL)
}
