package amqp

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const URI = "amqp://guest:guest@localhost:5672/"

// session is a composite of both connection and channel (inherited) this also forces a new object to have multiple channels.
type session struct {
	*amqp.Connection
	*amqp.Channel //todo bug here where if you change the queue and consume it doesn't work
}

// define a connection interface
type Connection interface {
	Publish(key, message string) error
	SetupQueue(queuename, key string) error
	Consume(queue string) error
}

// use a struct to isolate the changes
type internalAMQP struct {
	session  session
	exchange string
}

//Setup connections
func NewConnection(exchange string) Connection {

	internal := &internalAMQP{}
	internal.exchange = exchange

	return internal
}

func (am *internalAMQP) init() error {
	if am.session.Connection == nil || am.session.Connection.IsClosed() {
		conn, err := amqp.Dial(URI)
		if err != nil {
			return errors.Wrap(err, "amqp.init")
		}
		am.session.Connection = conn
	}

	channel, err := am.session.Connection.Channel()
	if err != nil {
		errors.Wrap(err, "amqp.init()")
	}

	am.session.Channel = channel

	err = am.session.Channel.ExchangeDeclare(am.exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "amqp.init()")
	}

	return nil

}

//Publish to the routing key a message
func (am *internalAMQP) Publish(key, message string) error {
	if err := am.init(); err != nil {
		return err
	}
	err := am.session.Channel.Publish(am.exchange, key, false, false, amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "application/json",
		ContentEncoding: "",
		Body:            []byte(message),
		DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
		Priority:        0,               // 0-9
		// a bunch of application/implementation-specific fields
	})

	if err != nil {
		return errors.Wrap(err, "amqp.Publish()")
	}

	fmt.Printf("Published Data: %s /w key: [%s]\n", message, key)

	return nil
}

//Setup the bindings
func (am *internalAMQP) SetupQueue(name, key string) error {
	if err := am.init(); err != nil {
		return err
	}

	_, err := am.session.Channel.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		errors.Wrap(err, "amqp.SetupQueue")
	}

	if err := am.session.Channel.QueueBind(name, key, am.exchange, true, nil); err != nil {
		return errors.Wrap(err, "amqp.SetupQueue")
	}

	return nil
}

//Consume the queue
func (am *internalAMQP) Consume(queue string) error {
	if err := am.init(); err != nil {
		return err
	}

	deliveries, err := am.session.Channel.Consume(queue, "dvp.consumer", false, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "amqp.Consume()")
	}

	for d := range deliveries {
		fmt.Printf(
			"QUEUE: %s got %dB delivery: [%v] %q\n",
			queue,
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		d.Ack(false)
	}

	return nil
}

// https://github.com/streadway/amqp/blob/master/_examples/pubsub/pubsub.go
/*

	if err := ch.ExchangeDeclare(
		"ORGExchange",     //exchange name
		"topic",           // exchange kind
		true,              // durable
		false,             // auto-delete
		false,             // internal
		false,             // asyncronous
		nil); err != nil { // map of extra options
		log.Fatalf("cannot declare fanout exchange: %v", err)
	}

*/
func NewSession() *session {
	ret := &session{}
	return ret
}

//s.Connect(); s.Publish(Exchange, message, key); s.Consume(Exchange, keyoattern)
func (s *session) Connect(url string) error {

	conn, err := amqp.Dial(url)
	if err != nil {
		return errors.Wrap(err, "amqp.Connect - ")
	}

	s.Connection = conn
	return nil
}
