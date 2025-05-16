package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeliveryLog struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey"`
	IncomingWebhookID uuid.UUID `gorm:"not null"`
	AttemptNumber     int       `gorm:"not null"`
	Status            string    `gorm:"not null"`
	HTTPStatusCode    *int
	ErrorDetails      string
	CreatedAt         time.Time `gorm:"index"`
}

func (dl *DeliveryLog) BeforeCreate(tx *gorm.DB) error {
	dl.ID = uuid.New()
	return nil
}
