package models

import "time"

type Post struct {
	ID        uint   `json:"id"`
	Title     string `gorm:"not null" json:"title"`
	Content   string `gorm:"not null" json:"content"`
	UserID    uint   `gorm:"not null" json:"user_id"`
	User      User
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
