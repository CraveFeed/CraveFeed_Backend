package Redis

import (
	"github.com/redis/go-redis/v9"
	"os"
)

var client *redis.Client
var connectionString = os.Getenv("REDIS_URL")

func GetClient() *redis.Client {
	if client == nil {
		opt, err := redis.ParseURL(connectionString)
		if err != nil {
			panic(err)
		}
		client = redis.NewClient(opt)
	}
	return client
}

func CloseClient() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
