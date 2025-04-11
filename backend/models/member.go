package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Member struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Position  string         `gorm:"type:varchar(255);not null" json:"position"`
	ImageURL  string         `gorm:"type:text" json:"imageUrl"`
	JoinedAt  time.Time      `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"joinedAt"`
	CreatedAt time.Time      `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (m *Member) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
