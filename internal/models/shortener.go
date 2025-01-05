package models

import (
	"errors"
	"regexp"
	"time"

	"github.com/bigxxby/dream-test-task/internal/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShortLink struct {
	ID        *uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID    *uuid.UUID `json:"user_id,omitempty" gorm:"type:uuid"`
	LongLink  string     `json:"long_url" gorm:"type:text;not null"`
	ShortId   string     `json:"short_id" gorm:"size:16;unique;not null"`
	Clicks    int        `json:"clicks" gorm:"default:0"`
	LastClick *time.Time `json:"last_click"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

func (u *ShortLink) BeforeCreate(tx *gorm.DB) (err error) {
	new := uuid.New()
	u.ID = &new
	return
}

func (u *ShortLink) UpdateClicks() {
	u.Clicks++
	now := time.Now()
	u.LastClick = &now
}

// regexp
func (u *ShortLink) ValidateLongLink() error {
	var regex = `^(https?|ftp)://[^\s/$.?#].[^\s]*$`
	match, _ := regexp.MatchString(regex, u.LongLink)
	if !match {
		return errors.New("invalid URL")
	}
	return nil

}

// adds app port and host to short link
func (u *ShortLink) ParseShortId() error {
	u.ShortId = "http://localhost:" + config.AppPort + "/shortener/" + u.ShortId
	return nil
}
