package delivery

import (
	"github.com/dathan/rabbit-topic/pkg/amqp"
	"github.com/pkg/errors"
)

type MailBox struct {
	QueueName  string
	RoutingKey string
}

type Delivery interface {
	SendTo(key string, body string) error
	MailBox(queueName string, key string) (*MailBox, error)
	Consume(mailBox *MailBox) error
}

type localDelivery struct {
	exchange   string
	connection amqp.Connection
}

//return a new Delivery Interface
func New(exchangeName string) Delivery {

	ret := &localDelivery{
		exchange:   exchangeName,
		connection: amqp.NewConnection(exchangeName),
	}
	return ret
}

//send the delivery
func (r *localDelivery) SendTo(queueName string, body string) error {
	if r.connection == nil {
		return errors.Errorf("connection was not set")
	}
	return r.connection.Publish(queueName, body)
}

//create a location for the delivery
func (r *localDelivery) MailBox(queue string, key string) (*MailBox, error) {
	box := &MailBox{
		queue,
		key,
	}
	if err := r.connection.SetupQueue(box.QueueName, box.RoutingKey); err != nil {
		return nil, err
	}

	return box, nil
}

//Consume the mailbox setup for delivery
func (r *localDelivery) Consume(box *MailBox) error {
	if r.connection == nil {
		return errors.New("Did not setup queue")
	}
	if err := r.connection.Consume(box.QueueName); err != nil {
		return err
	}
	return nil

}
