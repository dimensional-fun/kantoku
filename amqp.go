package main

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type AMQP struct {
	conn        *amqp091.Connection
	channel     *amqp091.Channel
	rpcQueue    amqp091.Queue
	rpcConsumer <-chan amqp091.Delivery

	Group   string
	Timeout time.Duration
}

var ErrNoRes = errors.New("no response from server")
var ErrDisconnected = errors.New("disconnected from the broker")

func (a *AMQP) Init(conn *amqp091.Connection) error {
	a.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	a.channel = ch

	return a.setupRPC()
}

func (a *AMQP) setupRPC() error {
	/* setup exchange. */
	err := a.channel.ExchangeDeclare(
		a.Group,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	/* setup queue */
	rpcQueue, err := a.channel.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	a.rpcQueue = rpcQueue

	/* setup consumer */
	rpcConsumer, err := a.channel.Consume(
		rpcQueue.Name,
		"",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	a.rpcConsumer = rpcConsumer
	return nil
}

func (a *AMQP) Call(event string, opts amqp091.Publishing) ([]byte, error) {
	if a.channel == nil {
		return nil, ErrDisconnected
	}

	correlation := uuid.New().String()

	opts.CorrelationId = correlation
	opts.ReplyTo = a.rpcQueue.Name
	opts.ContentType = "application/msgpack"
	opts.Expiration = "3000"

	err := a.channel.Publish(
		a.Group,
		event,
		false,
		true,
		opts,
	)

	if err != nil {
		return nil, err
	}

	for d := range a.rpcConsumer {
		if correlation == d.CorrelationId {
			return d.Body, nil
		}
	}

	return nil, ErrNoRes
}
