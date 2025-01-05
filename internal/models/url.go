package models

import (
	"time"

	"github.com/google/uuid"
)

// URL represents a short URL entity
type URL struct {
	ID        *uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID    *uuid.UUID `json:"user_id,omitempty" gorm:"type:uuid"`
	LongURL   string     `json:"long_url" gorm:"type:text;not null"`
	ShortID   string     `json:"short_id" gorm:"size:16;unique;not null"`
	Clicks    int        `json:"clicks" gorm:"default:0"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

func (u *URL) BeforeCreate() (err error) {
	new := uuid.New()
	u.ID = &new
	return
}
