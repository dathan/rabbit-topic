package main

import (
	"github.com/dathan/rabbit-topic/pkg/delivery"
)

// need two objects or two connections since the channels are tied to the connection, there is a bug if you use the same object across queue names.
func main() {

	waitfor1 := delivery.New("wework.topic.exchange")
	box1, err := waitfor1.MailBox("ct-131.production.queue", "wework.production.*")

	if err != nil {
		panic(err)
	}

	waitfor2 := delivery.New("wework.topic.exchange")
	box2, err := waitfor2.MailBox("ct-131.dathan-development.queue", "wework.production.*")
	if err != nil {
		panic(err)
	}

	_, err = waitfor2.MailBox("ct-131.dathan-development.queue", "wework.NOTproduction.user")
	go waitfor1.Consume(box1)
	waitfor2.Consume(box2)

}
