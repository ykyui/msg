package redis

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

func Init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDISHOST"), os.Getenv("REDISPORT")),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if _, err := rdb.Ping().Result(); err != nil {
		panic(err)
	} else {
		fmt.Println("redis ready")
	}
	go userOnlineService()
}

func Close() {
	rdb.Close()
}
