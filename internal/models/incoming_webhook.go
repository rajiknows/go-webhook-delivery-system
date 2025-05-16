package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IncomingWebhook struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	SubscriptionID uuid.UUID `gorm:"not null"`
	EventType      string
	Payload        string `gorm:"type:jsonb"`
	CreatedAt      time.Time
}

func (iw *IncomingWebhook) BeforeCreate(tx *gorm.DB) error {
	iw.ID = uuid.New()
	return nil
}
