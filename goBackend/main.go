package main

import (
	"log"
	"msg/controller"
	"msg/mongodb"
	"msg/mq"
	"msg/redis"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	mq.Init()
	defer mq.Close()
	mongodb.Init()
	defer mongodb.Close()
	redis.Init()
	defer redis.Close()

	controller.Run()

	t := make(chan struct{}, 1)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		t <- struct{}{}
	}()
	<-t
}
