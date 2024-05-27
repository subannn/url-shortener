package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	redis "github.com/subannn/urlshorter/redis"
	server "github.com/subannn/urlshorter/server"
)

func main() {
	redis.RunRedis()
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
