package factory

import "github.com/redis/go-redis/v9"

func CreateNewRedisClient(hostName string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: hostName,
		DB:   0,
	})
}
