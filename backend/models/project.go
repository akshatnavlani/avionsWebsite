package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text;not null" json:"description"`
	MarkdownURL string         `gorm:"type:text" json:"markdownUrl"`
	ImageURL    string         `gorm:"type:text" json:"imageUrl"`
	CreatedAt   time.Time      `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
