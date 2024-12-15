package rabbit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/model"
	"github.com/streadway/amqp"
)

type Rabbit struct {
	nameExchange string
	typeExchange string
	nameQueue    string
	routingKey   string
	uri          string
	logg         *logger.Logger
	conn         *amqp.Connection
}

func New(config *config.Config, logg *logger.Logger) *Rabbit {
	rmq := config.Rabbit
	return &Rabbit{
		nameExchange: config.Rabbit.Name,
		typeExchange: "direct",
		nameQueue:    config.Rabbit.Name + "-queue",
		routingKey:   config.Rabbit.Name + "-key",
		uri:          fmt.Sprintf("amqp://%s:%s@%s:%d/", config.Rabbit.User, rmq.Password, rmq.Host, rmq.Port),
		logg:         logg,
	}
}

func (r *Rabbit) Notify(ctx context.Context, events []*model.Event) error {
	defer ctx.Done()

	channel, err := r.conn.Channel()
	if err != nil {
		r.logg.Error("rabbit notify channel", err)
		return nil
	}
	defer channel.Close()

	if err := r.exchangeDeclare(channel); err != nil {
		r.logg.Error("rabbit notify exchange declare", err)
		return nil
	}

	for _, event := range events {
		body, err := json.Marshal(model.GetNotificationFromEvent(event))
		if err != nil {
			r.logg.Error("rabbit notify marshal event", err)
			continue
		}
		if err = channel.Publish(
			r.nameExchange, // publish to an exchange
			r.routingKey,   // routing to 0 or more queues
			false,          // mandatory
			false,          // immediate
			amqp.Publishing{
				Headers:         amqp.Table{},
				ContentType:     "text/plain",
				ContentEncoding: "",
				Body:            body,
				DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
				Priority:        0,              // 0-9
			},
		); err != nil {
			r.logg.Error("rabbit notify channel publish", err)
			return nil
		}
	}

	return nil
}

func (r *Rabbit) Read(ctx context.Context) error {
	channel, err := r.conn.Channel()
	if err != nil {
		r.logg.Error("rabbit read channel", err)
		return nil
	}

	go func() {
		<-ctx.Done()
		if err := channel.Close(); err != nil {
			r.logg.Error("rabbit read channel close", err)
		}
	}()

	if err := r.exchangeDeclare(channel); err != nil {
		r.logg.Error("rabbit read exchange declare", err)
		return nil
	}

	queue, err := r.queueDeclare(channel)
	if err != nil {
		r.logg.Error("rabbit read queue declare", err)
		return nil
	}

	if err := r.queueBind(channel); err != nil {
		r.logg.Error("rabbit read queue bind", err)
		return nil
	}

	messages, err := channel.Consume(queue.Name, r.nameExchange, false, false, false, false, nil)
	if err != nil {
		r.logg.Error("rabbit read queue consume", err)
		return nil
	}

	r.handleMessages(messages)

	<-ctx.Done()
	return nil
}

func (r *Rabbit) handleMessages(messages <-chan amqp.Delivery) {
	go func() {
		for message := range messages {
			notification := &model.Notification{}
			if err := json.Unmarshal(message.Body, notification); err != nil {
				r.logg.Error("sender send", err)
			}

			r.logg.Info(fmt.Sprintf("Notification sent. Id: %s; Title: %s; Start: %s; UserId: %s;",
				notification.ID,
				notification.Title,
				notification.Start,
				notification.UserID))

			message.Ack(false)
		}
	}()
}

func (r *Rabbit) exchangeDeclare(channel *amqp.Channel) error {
	return channel.ExchangeDeclare(
		r.nameExchange, // name
		r.typeExchange, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // noWait
		nil,            // arguments
	)
}

func (r *Rabbit) queueDeclare(channel *amqp.Channel) (amqp.Queue, error) {
	return channel.QueueDeclare(
		r.nameQueue, // name of the queue
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // noWait
		nil,         // arguments
	)
}

func (r *Rabbit) queueBind(channel *amqp.Channel) error {
	return channel.QueueBind(
		r.nameQueue,    // name of the queue
		r.routingKey,   // bindingKey
		r.nameExchange, // sourceExchange
		false,          // noWait
		nil,            // arguments
	)
}

func (r *Rabbit) Connect() (err error) {
	r.conn, err = amqp.Dial(r.uri)
	if err != nil {
		r.logg.Error("Failed to connect rabbit: ", err)
		return
	}
	r.logg.Info("Rabbit connect")
	return
}

func (r *Rabbit) Close() error {
	r.logg.Info("Rabbit close connection")
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
