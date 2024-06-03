package redis

import (
	"context"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	models "github.com/subannn/urlshorter/models"
)

var ctx = context.Background()
var rdb *redis.Client

var hashName = os.Getenv("HASH_NAME")
var sortedSetName = os.Getenv("SORTED_SET_NAME")

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

func CutAndSaveURL(longURL models.RequestLongURL) string {
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

	err = rdb.HSet(ctx, hashName, shortURL, longURL.LongURL).Err()
	if err != nil {
		panic(err)
	}

	return shortURL
}

func SaveExpirationDate(shortURL string, hrs int) {
	z := redis.Z{
		Score:  float64(hrs),
		Member: shortURL,
	}

	err := rdb.ZAdd(ctx, sortedSetName, z).Err()
	if err != nil {
		log.Panic(err)	
	}
}

func ShutdownRedis(ctxToShutdown context.Context) {
	rdb.Shutdown(ctxToShutdown)
	log.Println("Redis closed")
}
