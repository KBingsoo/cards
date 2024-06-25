package cards

import (
	"github.com/KBingsoo/cards/pkg/models/event"
	"github.com/KBingsoo/entities/pkg/models"
	"github.com/literalog/go-wise/wise"
)

type Repository wise.MongoRepository[models.Card]

type Producer interface {
	Emit(event event.Event) error
}
