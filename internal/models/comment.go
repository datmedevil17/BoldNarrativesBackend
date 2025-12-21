package models

import (
    "time"
    "gorm.io/gorm"
)

type Comment struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Comment   string         `json:"comment" gorm:"type:text;not null"`
    AuthorID  uint           `json:"author_id" gorm:"not null;index"`
    BlogID    uint           `json:"blog_id" gorm:"not null;index"`
    Author    User           `json:"author" gorm:"foreignKey:AuthorID"`
    Blog      Blog           `json:"-" gorm:"foreignKey:BlogID"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CommentResponse struct {
    ID        uint         `json:"id"`
    Comment   string       `json:"comment"`
    AuthorID  uint         `json:"author_id"`
    Author    UserResponse `json:"author"`
    CreatedAt time.Time    `json:"created_at"`
}