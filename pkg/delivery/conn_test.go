package delivery

import (
	"testing"

	"github.com/dathan/rabbit-topic/pkg/amqp"
)

func Test_localDelivery_SendTo(t *testing.T) {
	type fields struct {
		exchange   string
		connection amqp.Connection
	}
	type args struct {
		key  string
		body string
	}
	str := `{"page": 1, "fruits": ["apple", "peach"]}`
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"PUBLISH", fields{"wework.topic.exchange", amqp.NewConnection("wework.topic.exchange")}, args{"wework.dathanenviron.users", str}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &localDelivery{
				exchange:   tt.fields.exchange,
				connection: tt.fields.connection,
			}
			if err := r.SendTo(tt.args.key, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("localDelivery.SendTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_localDelivery_Consume(t *testing.T) {
	type fields struct {
		exchange   string
		connection amqp.Connection
	}
	type args struct {
		box *MailBox
	}

	dev := New("wework.topic.exchange")
	box, err := dev.MailBox("dathan.awesome.queue", "wework.*")
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"CONSUME", fields{"wework.topic.exchange", amqp.NewConnection("wework.topic.exchange")}, args{box}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &localDelivery{
				exchange:   tt.fields.exchange,
				connection: tt.fields.connection,
			}
			if err := r.Consume(tt.args.box); (err != nil) != tt.wantErr {
				t.Errorf("localDelivery.Consume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
