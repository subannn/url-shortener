package main

import (
	redis "github.com/subannn/urlshorter/redis"
	server "github.com/subannn/urlshorter/server"

)

func main() {
	redis.RunRedis()
	server.RunServer()
}
