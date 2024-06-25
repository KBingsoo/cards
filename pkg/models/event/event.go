package event

import (
	"context"
	"time"

	"github.com/KBingsoo/entities/pkg/models"
)

type EventType string

const (
	Create       EventType = "card_create"
	Update       EventType = "card_update"
	OrderFulfill EventType = "card_order_fulfill"
	Delete       EventType = "card_delete"

	Succeed EventType = "card_succeed"
	Error   EventType = "card_error"
)

type Event struct {
	Type    EventType       `json:"type"`
	Time    time.Time       `json:"time"`
	Card    models.Card     `json:"card"`
	Context context.Context `json:"context"`
}
