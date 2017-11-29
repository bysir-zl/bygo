package mq

import (
	"testing"
	"log"
	"github.com/streadway/amqp"
)

func TestPublish(t *testing.T) {
	r := NewRabbit("amqp://guest:guest@localhost:5672/")
	r.Publish("hello", []byte("world"))
}

func TestReceive(t *testing.T) {
	r := NewRabbit("amqp://guest:guest@localhost:5672/")
	err := r.Receive("hello", func(delivery amqp.Delivery) error {
		log.Print(string(delivery.Body))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	select {}
}
