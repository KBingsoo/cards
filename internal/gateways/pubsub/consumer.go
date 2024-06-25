package pubsub

import (
	"github.com/KBingsoo/cards/internal/domain/cards"
	"github.com/KBingsoo/cards/pkg/models/event"
	"github.com/seosoojin/go-rabbit/rabbit"
	rabbitConsumer "github.com/seosoojin/go-rabbit/rabbit/consumer"
	"github.com/streadway/amqp"
)

type consumer struct {
	internalConsumer rabbit.Consumer[event.Event]
	manager          cards.Manager
}

func NewConsumer(conn *amqp.Connection) (*consumer, error) {
	internalConsumer, err := rabbit.NewConsumer[event.Event](conn, decode, rabbitConsumer.WithConsumerExchange("cards", "topic", ""))
	if err != nil {
		return nil, err
	}

	return &consumer{
		internalConsumer: internalConsumer,
	}, nil
}

func (c *consumer) Consume() error {
	return c.internalConsumer.Consume(c.handler)
}

func (c *consumer) handler(entry event.Event) error {
	switch entry.Type {
	case event.OrderFulfill:
		card, err := c.manager.GetByID(entry.Context, entry.Card.ID)
		if err != nil {
			return err
		}

		card.Qty--

		return c.manager.Update(entry.Context, &card)

	default:
		return nil
	}

}
