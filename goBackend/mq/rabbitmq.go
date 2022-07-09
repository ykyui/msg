package mq

import (
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
)

func Init() {
	if _conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", os.Getenv("MQUSERID"), os.Getenv("MQPASSWORD"), os.Getenv("MQHOST"), os.Getenv("MQPORT"))); err != nil {
		panic(err)
	} else {
		conn = _conn
	}
	if _ch, err := conn.Channel(); err != nil {
		panic(err)
	} else {
		ch = _ch
	}
	if err := ch.ExchangeDeclare("msg", "headers", true, false, false, false, nil); err != nil {
		panic(err)
	}
	if err := ch.ExchangeDeclare("online", "topic", true, false, false, false, nil); err != nil {
		panic(err)
	}
	fmt.Println("mq ready")
}

func CatchOnlineStatusConsumer(queueName string) (<-chan amqp.Delivery, error) {
	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	if err = ch.QueueBind(q.Name, "online.#", "online", false, nil); err != nil {
		return nil, err
	}
	return ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}

func Close() {
	fmt.Println("close")
	// ch.ExchangeDelete("msg", true, true)
	// ch.ExchangeDelete("online", true, true)
	ch.Close()
	conn.Close()
}

func UserOnline(username string) error {
	return ch.Publish("online", "online."+username, false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte("online")})
}

func UserOffline(username string) error {
	return ch.Publish("online", "online."+username, false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte("offline")})
}
