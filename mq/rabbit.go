package mq

import (
	"github.com/streadway/amqp"
	"time"
	"log"
)

type Rabbit struct {
	url         string
	publishConn *amqp.Connection
}

// 发送也是长连接，短连接在连接上非常耗时
func (p *Rabbit) Publish(queue string, body []byte) (error) {
	var err error
	if p.publishConn == nil {
		p.publishConn, err = amqp.Dial(p.url)
	}
	if err != nil {
		return err
	}
	//defer conn.Close()
	ch, err := p.publishConn.Channel()
	if err != nil {
		return err
	}
	q, err := ch.QueueDeclare(
		queue, // name
		true, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        body,
	})
	if err != nil {
		return err
	}
	return nil
}

type Handler func(amqp.Delivery) error

// 接收消息, 如果断线会在5s后重试
func (p *Rabbit) Receive(queue string, h Handler) (error) {
	conn, err := amqp.Dial(p.url)
	if err != nil {
		return err
	}
	conn.Close()

	go func() {
		for {
			conn, err := amqp.Dial(p.url)
			if err != nil {
				log.Printf("[rabbit] receive error: %v, try again after 5s", err)
				time.Sleep(5 * time.Second)
				continue
			}
			ch, err := conn.Channel()
			if err != nil {
				log.Printf("[rabbit] receive error: %v, try again after 5s", err)
				time.Sleep(5 * time.Second)
				continue
			}
			q, err := ch.QueueDeclare(
				queue, // name
				true, // durable
				false, // delete when unused
				false, // exclusive
				false, // no-wait
				nil,   // arguments
			)
			if err != nil {
				log.Printf("[rabbit] receive error: %v, try again after 5s", err)
				time.Sleep(5 * time.Second)
				continue
			}
			msg, err := ch.Consume(q.Name, "", true, false, false, false, nil)
			if err != nil {
				log.Printf("[rabbit] receive error: %v, try again after 5s", err)
				time.Sleep(5 * time.Second)
				continue
			}
			for i := range msg {
				err := h(i)
				if err != nil {
					log.Printf("[rebit] woker respon is error: %v", err)
				}
			}

			time.Sleep(5 * time.Second)
		}
	}()

	return nil
}

func NewRabbit(url string) *Rabbit {
	return &Rabbit{url: url}
}
