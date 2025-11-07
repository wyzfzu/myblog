package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Title     string         `gorm:"size:64;not null" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	UserID    uint           `json:"userId"`
	User      User           `json:"author,omitempty"`
	Comments  []Comment      `json:"comments,omitempty"`
}

type CreatePostReq struct {
	Title   string `json:"title" binding:"required,min=2,max=64"`
	Content string `json:"content" binding:"required,min=2"`
}

type UpdatePostReq struct {
	ID      uint   `json:"id" binding:"required"`
	Title   string `json:"title" binding:"required,min=2,max=64"`
	Content string `json:"content" binding:"required,min=2"`
}
