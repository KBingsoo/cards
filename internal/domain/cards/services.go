package cards

import (
	"context"
	"time"

	"github.com/KBingsoo/cards/pkg/models/event"
	"github.com/KBingsoo/entities/pkg/models"
	"github.com/google/uuid"
)

type Manager interface {
	Create(ctx context.Context, card *models.Card) error
	GetAll(ctx context.Context) ([]models.Card, error)
	GetByID(ctx context.Context, id string) (models.Card, error)
	Update(ctx context.Context, card *models.Card) error
	Delete(ctx context.Context, id string) (models.Card, error)
}

type manager struct {
	repository Repository
	producer   Producer
}

func NewManager(repository Repository, producer Producer) *manager {
	return &manager{
		repository: repository,
		producer:   producer,
	}
}

func (m *manager) Create(ctx context.Context, card *models.Card) error {
	if card.ID == "" {
		card.ID = uuid.NewString()
	}

	return m.repository.Upsert(ctx, card.ID, *card)
}

func (m *manager) GetAll(ctx context.Context) ([]models.Card, error) {
	return m.repository.FindAll(ctx)
}

func (m *manager) GetByID(ctx context.Context, id string) (models.Card, error) {
	return m.repository.Find(ctx, id)
}

func (m *manager) Update(ctx context.Context, card *models.Card) error {
	msg := event.Event{
		Type: event.Succeed,
		Time: time.Now(),
		Card: *card,
	}

	err := m.repository.Upsert(ctx, card.ID, *card)
	if err != nil {
		msg.Type = event.Error
	}

	msg.Context = ctx

	return m.producer.Emit(msg)
}

func (m *manager) Delete(ctx context.Context, id string) (models.Card, error) {
	return m.repository.Delete(ctx, id)
}
