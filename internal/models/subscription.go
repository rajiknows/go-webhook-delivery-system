package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	TargetURL       string    `gorm:"not null"`
	Secret          string
	Active          bool
	EventTypeFilter string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.New()
	return nil
}
