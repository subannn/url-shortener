package redis

import (
	"context"
	"os"

	"github.com/google/uuid"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func RunRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})
}

func GetLongURL(shortURL string) string {
	val, err := rdb.HGet(ctx, "URLS", shortURL).Result()
	if err != nil {
		panic(err)
	}

	return val

}

func CutAndSaveURL(longURL string) string {

	id := uuid.New()
	shortURL := string(id.String()[1:6])
	exists, err := rdb.HExists(ctx, "URLS", shortURL).Result()
	if err != nil {
		panic(err)
	}

	for exists {
		id = uuid.New()
		shortURL = string(id.String()[1:6])

		exists, err = rdb.HExists(ctx, "URLS", shortURL).Result()
		if err != nil {
			panic(err)
		}
	}

	err = rdb.HSet(ctx, "URLS", shortURL, longURL).Err()
	if err != nil {
		panic(err)
	}

	return shortURL
}
