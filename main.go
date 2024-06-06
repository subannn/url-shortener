package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	redis "github.com/subannn/urlshorter/redis"
	server "github.com/subannn/urlshorter/server"
)

func main() {
	fmt.Println(os.Getenv("REDIS_ADDRESS"))
	mtx := &sync.Mutex{}
	redis.RunRedis(mtx)

	go func() {
		redis.DeleteExpitedURLS(int(time.Now().Unix() / 3600))
		for {
			ticker := time.Tick(time.Hour)
			<-ticker
			redis.DeleteExpitedURLS(int(time.Now().Unix() / 3600))
		}
	}()

	go server.RunServer()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server.ShutDownServer(ctx)
	redis.ShutdownRedis(ctx)

	<-ctx.Done()
	log.Println("shutting down")
}
