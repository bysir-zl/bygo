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

// 测试群发

func TestPublishFull(t *testing.T) {
	r := NewRabbitFull("amqp://guest:guest@localhost:5672/")
	ch, err := r.NewChannel()
	if err != nil {
		t.Fatal(err)
	}

	err = ch.ExchangeDeclare("msg_box", "fanout", false).
		Publish([]byte("i am is rabbit publisher"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestReceiveFull1(t *testing.T) {
	r := NewRabbitFull("amqp://guest:guest@localhost:5672/")
	r.SimplyReceive("msg_box", "fanout", "1", func(delivery amqp.Delivery) error {
		log.Print(string(delivery.Body))
		return nil
	})
	select {}
}

func TestReceiveFull2(t *testing.T) {
	r := NewRabbitFull("amqp://guest:guest@localhost:5672/")
	r.SimplyReceive("msg_box", "fanout", "2", func(delivery amqp.Delivery) error {
		log.Print(string(delivery.Body))
		return nil
	})
	select {}
}
