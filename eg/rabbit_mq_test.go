package eg

import (
	"github.com/streadway/amqp"
	"testing"
	"log"
)

func TestPush(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		t.Fatal(err)
	}

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		t.Fatal(err)
	}
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("hello"),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("ok")

}

func TestReceive(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		t.Fatal(err)
	}

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	for i := range msgs {
		log.Printf("Received a message: %s", i.Body)
	}
}
