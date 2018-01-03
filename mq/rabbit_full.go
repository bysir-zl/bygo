package mq

import (
	"github.com/streadway/amqp"
	"time"
	"log"
	"errors"
)

type RabbitFull struct {
	url         string
	publishConn *amqp.Connection
}

type Channel struct {
	*amqp.Channel
	err      error // 在链式操作中的err
	exchange string
	queue    string
}

// 连接
func (p *RabbitFull) NewChannel() (*Channel, error) {
	var err error
	if p.publishConn == nil {
		p.publishConn, err = amqp.Dial(p.url)
	}
	if err != nil {
		return nil, err
	}

	ch, err := p.publishConn.Channel()
	if err != nil {
		p.publishConn.Close()
		p.publishConn = nil
		return nil, err
	}
	return &Channel{Channel: ch}, nil
}

func (ch *Channel) QueueDeclare(queue string, durable bool) (*Channel) {
	_, err := ch.Channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.err = errors.New("QueueDeclare err:" + err.Error())
	} else {
		ch.queue = queue
	}
	return ch
}

func (ch *Channel) ExchangeDeclare(exchange string, kind string, durable bool) (*Channel) {
	err := ch.Channel.ExchangeDeclare(exchange, kind, durable, false, false, false, nil)
	if err != nil {
		ch.err = errors.New("ExchangeDeclare err:" + err.Error())
	} else {
		ch.exchange = exchange
	}
	return ch
}

// bind刚刚声明的queue和exchange
func (ch *Channel) Bind() (*Channel) {
	err := ch.Channel.QueueBind(ch.queue, "", ch.exchange, false, nil)
	if err != nil {
		ch.err = errors.New("bind err:" + err.Error())
	}

	return ch
}

func (ch *Channel) Publish(body []byte) (error) {
	if ch.err != nil {
		return ch.err
	}
	err := ch.Channel.Publish(ch.exchange, ch.queue, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        body,
	})
	if err != nil {
		return err
	}

	//ch.Close()
	return nil
}

func (ch *Channel) Receive(h Handler) (error) {
	if ch.err != nil {
		return ch.err
	}
	msg, err := ch.Consume(ch.queue, "", true, false, false, false, nil)
	if err != nil {
		return nil
	}
	for i := range msg {
		err := h(i)
		if err != nil {
			log.Printf("[rebit] woker response is error: %v", err)
		}
	}

	return nil
}

// 发送也是长连接，短连接在连接上非常耗时
func (p *RabbitFull) SimplyPublish(queue string, body []byte) (error) {
	ch, err := p.NewChannel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.QueueDeclare(queue, false).Publish(body)
}

// 接收消息, 如果断线会在5s后重试
func (p *RabbitFull) SimplyReceive(exchange string, exchangeKind string, queue string, h Handler) (error) {
	go func() {
		for {
			ch, err := p.NewChannel()
			if err != nil {
				log.Printf("[rabbit] receive error: %v, try again after 5s", err)
				time.Sleep(5 * time.Second)
				continue
			}
			if queue != "" {
				ch = ch.QueueDeclare(queue, false)
			}
			if exchange != "" {
				ch = ch.ExchangeDeclare(exchange, exchangeKind, false)
			}
			if exchange != "" && queue != "" {
				ch = ch.Bind()
			}

			err = ch.Receive(h)
			if err != nil {
				log.Printf("[rabbit] receive error: %v, try again after 5s", err)
			}

			ch.Close()
			time.Sleep(5 * time.Second)
		}
	}()

	return nil
}

func NewRabbitFull(url string) *RabbitFull {
	return &RabbitFull{url: url}
}
