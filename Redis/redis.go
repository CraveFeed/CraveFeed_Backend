package Redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func GetClient() *redis.Client {
	if client == nil {
		fmt.Println("Creating a redis Client")
		client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	}
	return client
}

func CloseClient() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
