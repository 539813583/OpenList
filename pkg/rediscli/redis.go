package rediscli

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func Init(addr, password string, db int) {
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func Client() *redis.Client {
	if client == nil {
		fmt.Println("warning: redis client not initialized, webdav session control disabled")
	}
	return client
}
