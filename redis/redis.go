package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

var hashName = os.Getenv("HASH_NAME")

func RunRedis() {
	rdb = redis.NewClient(&redis.Options{

		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})
}

func GetLongURL(shortURL string) string {
	val, err := rdb.HGet(ctx, hashName, shortURL).Result()
	if err != nil {
		panic(err)
	}

	return val

}

func CutAndSaveURL(longURL string) string {

	id := uuid.New()
	shortURL := string(id.String()[1:6])
	exists, err := rdb.HExists(ctx, hashName, shortURL).Result()
	if err != nil {
		panic(err)
	}

	for exists { // hash can be used
		id = uuid.New()
		shortURL = string(id.String()[1:6])

		exists, err = rdb.HExists(ctx, hashName, shortURL).Result()
		if err != nil {
			panic(err)
		}
	}

	err = rdb.HSet(ctx, hashName, shortURL, longURL).Err()
	if err != nil {
		fmt.Println("DOWS NOOOT WOORKKK")
		panic(errors.New("DOWS NOOOT WOORKKK"))
	}

	return shortURL
}

func ShutdownRedis(ctxToShutdown context.Context) {
	rdb.Shutdown(ctxToShutdown)
	log.Println("Redis closed")
}
