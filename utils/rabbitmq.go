package utils

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

type RabbitMq struct {
	Channel  *amqp.Channel
	Name     string
	Exchange string
}

func NewRabbitMq(s string) *RabbitMq {

	conn, errDial := amqp.Dial(s)
	if errDial != nil {
		panic(errDial)
	}

	ch, errChan := conn.Channel()
	if errChan != nil {
		panic(errChan)
	}

	q, errQueueDeclare := ch.QueueDeclare(
		"",
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if errQueueDeclare != nil {
		panic(errQueueDeclare)
	}

	mq := new(RabbitMq)
	mq.Channel = ch
	mq.Name = q.Name
	return mq

}

func (q *RabbitMq) Bind(exchange string) {

	e := q.Channel.QueueBind(
		q.Name,   // queue Name
		"",       // routing key
		exchange, // exchange
		false,
		nil,
	)
	if e != nil {
		panic(e)
	}
	q.Exchange = exchange
}

func (q *RabbitMq) Send(queue string, body interface{}) {

	str, e := json.Marshal(body)

	if e != nil {
		panic(e)
	}

	e = q.Channel.Publish("",
		queue,
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	if e != nil {
		panic(e)
	}
}

func (q *RabbitMq) Publish(exchange string, body interface{}) {

	str, e := json.Marshal(body)

	if e != nil {
		panic(e)
	}

	e = q.Channel.Publish(exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})

	if e != nil {
		panic(e)
	}
}

func (q *RabbitMq) Consume() <-chan amqp.Delivery {
	c, e := q.Channel.Consume(q.Name, "",
		true,
		false,
		false,
		false,
		nil)

	if e != nil {
		panic(e)
	}

	return c

}

func (q *RabbitMq) Close() {
	q.Channel.Close()
}
