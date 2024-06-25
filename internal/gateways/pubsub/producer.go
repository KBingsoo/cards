package pubsub

import (
	"encoding/json"
	"fmt"

	"github.com/KBingsoo/cards/pkg/models/event"
	"github.com/seosoojin/go-rabbit/rabbit"
	"github.com/seosoojin/go-rabbit/rabbit/message"
	rabbitProducer "github.com/seosoojin/go-rabbit/rabbit/producer"
	"github.com/streadway/amqp"
)

type producer struct {
	internalProducer rabbit.Producer
}

func NewProducer(conn *amqp.Connection) (*producer, error) {
	internalProducer, err := rabbit.NewProducer(conn, rabbitProducer.WithProducerExchange("cards", "topic", ""))
	if err != nil {
		return nil, err
	}

	return &producer{
		internalProducer: internalProducer,
	}, nil
}

func (p *producer) Emit(event event.Event) error {

	b, err := json.Marshal(event)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("card.%s.%s.%s", event.Card.ID, event.Type, event.Time.Format("2006-01-02"))

	msg := message.Message{
		Key: key,
		Value: amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		},
	}

	return p.internalProducer.Emit(msg)
}
