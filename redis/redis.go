package redis

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	models "github.com/subannn/urlshorter/models"
)

var ctx = context.Background()
var rdb *redis.Client
var mtx *sync.Mutex

var hashName = os.Getenv("HASH_NAME")
var sortedSetName = os.Getenv("SORTED_SET_NAME")

func RunRedis(newMutex *sync.Mutex) {
	mtx = newMutex
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
	mtx.Lock()
	defer mtx.Unlock()

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
	mtx.Lock()
	defer mtx.Unlock()
	z := redis.Z{
		Score:  float64(hrs),
		Member: shortURL,
	}

	err := rdb.ZAdd(ctx, sortedSetName, z).Err()
	if err != nil {
		log.Panic(err)
	}
}

func DeleteExpitedURLS(hoursFromUnixTime int) {
	mtx.Lock()
	defer mtx.Unlock()
	z := &redis.ZRangeBy{
		Min: "0",
		Max: strconv.Itoa(hoursFromUnixTime),
	}

	ExpiredURLs, err := rdb.ZRangeByScore(ctx, sortedSetName, z).Result()
	if err != nil {
		log.Panic(err)
	}

	var ExpiredURLs_Interfaces []interface{}
	for _, url := range ExpiredURLs {
		ExpiredURLs_Interfaces = append(ExpiredURLs_Interfaces, url)
	}

	err = rdb.HDel(ctx, hashName, ExpiredURLs...).Err() // Delete URLs from hash set
	if err != nil {
		log.Panic(err)
	}

	err = rdb.ZRem(ctx, sortedSetName, ExpiredURLs_Interfaces...).Err() // Delete URLs from sorted set
	if err != nil {
		log.Panic(err)
	}
}

func ShutdownRedis(ctxToShutdown context.Context) {
	rdb.Shutdown(ctxToShutdown)
	log.Println("Redis closed")
}
