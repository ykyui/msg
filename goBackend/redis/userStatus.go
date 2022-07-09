package redis

import (
	"fmt"
	"msg/controller"
	"msg/mq"
	"strings"
	"time"
)

func userOnlineService() {
	msgs, err := mq.CatchOnlineStatusConsumer("online.user")
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		routing := strings.Split(msg.RoutingKey, ".")
		if string(msg.Body) != "online" {
			fmt.Println("offline")
			rdb.Del(routing[1])
			continue
		}
		rdb.Set(fmt.Sprintf("lastonline.%s", routing[1]), time.Now(), getExTime())
		if controller.ReverseProxy == 1 {
			go func() {
				time.Sleep(time.Duration(controller.LongpollTimeout+1) * time.Second)
				_, err := rdb.Get(fmt.Sprintf("lastonline.%s", routing[1])).Result()
				if err != nil {
					mq.UserOffline(routing[1])
				}
			}()
		}
	}
}

func getExTime() time.Duration {
	if controller.ReverseProxy == 0 {
		return 0
	} else {
		return time.Duration(controller.LongpollTimeout) * time.Second
	}
}
