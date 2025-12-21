package models

import (
    "time"
    "gorm.io/gorm"
)

type Blog struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Title     string         `json:"title" gorm:"not null;index"`
    Content   string         `json:"content" gorm:"type:text;not null"`
    Genre     string         `json:"genre" gorm:"not null;index"`
    Views     int            `json:"views" gorm:"default:0"`
    AuthorID  uint           `json:"author_id" gorm:"not null;index"`
    Author    User           `json:"author" gorm:"foreignKey:AuthorID"`
    Votes     []Vote         `json:"votes,omitempty" gorm:"foreignKey:BlogID;constraint:OnDelete:CASCADE"`
    Comments  []Comment      `json:"comments,omitempty" gorm:"foreignKey:BlogID;constraint:OnDelete:CASCADE"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type BlogResponse struct {
    ID        uint         `json:"id"`
    Title     string       `json:"title"`
    Content   string       `json:"content,omitempty"`
    Genre     string       `json:"genre"`
    Views     int          `json:"views"`
    VoteCount int64        `json:"votes"`
    AuthorID  uint         `json:"author_id"`
    Author    UserResponse `json:"author"`
    CreatedAt time.Time    `json:"created_at"`
}

type BlogListResponse struct {
    ID        uint         `json:"id"`
    Title     string       `json:"title"`
    Genre     string       `json:"genre"`
    Views     int          `json:"views"`
    VoteCount int64        `json:"votes"`
    AuthorID  uint         `json:"author_id"`
    Author    UserResponse `json:"author"`
    CreatedAt time.Time    `json:"created_at"`
}