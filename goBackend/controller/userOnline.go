package controller

import (
	"fmt"
	"msg/mq"
	"net/http"
	"time"
)

func init() {
	countInit.Add(1)
	go func() {
		addPrivateApi("/longpollMsg", longpollMsg, []string{http.MethodPost})
		countInit.Done()
	}()
}

func longpollMsg(rw http.ResponseWriter, r *http.Request, username string) (interface{}, int, error) {
	r.Body.Close()
	fmt.Println("userId ", username, " connect")
	mq.UserOnline(username)
	t := time.NewTicker(time.Duration(LongpollTimeout) * time.Second)
	for {
		select {
		case <-t.C:
			if ReverseProxy == 1 {
				return struct{}{}, http.StatusOK, nil
			}
		case <-r.Context().Done():
			mq.UserOffline(username)
			return "cancelled", http.StatusOK, nil
		}
	}
}
