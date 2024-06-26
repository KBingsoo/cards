package pubsub

import (
	"github.com/KBingsoo/cards/internal/domain/cards"
	"github.com/KBingsoo/cards/pkg/models/event"
	"github.com/KBingsoo/entities/pkg/models"
	"github.com/seosoojin/go-rabbit/rabbit"
	rabbitConsumer "github.com/seosoojin/go-rabbit/rabbit/consumer"
	"github.com/streadway/amqp"
)

type consumer struct {
	internalConsumer rabbit.Consumer[event.Event]
	manager          cards.Manager
	orderItems       map[string][]models.Card
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

		if _, ok := c.orderItems[entry.OrderID]; !ok {
			c.orderItems[entry.OrderID] = make([]models.Card, 0)
		}

		c.orderItems[entry.OrderID] = append(c.orderItems[entry.OrderID], card)

		card.Qty--

		c.manager.Update(entry.Context, &card)

		return nil

	case event.OrderRevert:
		for _, card := range c.orderItems[entry.OrderID] {
			err := c.manager.Update(entry.Context, &card)
			if err != nil {
				return err
			}
		}

		delete(c.orderItems, entry.OrderID)

		return nil

	default:
		return nil
	}

}
