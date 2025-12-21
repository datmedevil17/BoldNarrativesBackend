package models

import (
    "time"
    "gorm.io/gorm"
)

type Vote struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    UserID    uint           `json:"user_id" gorm:"not null;index:idx_user_blog,unique"`
    BlogID    uint           `json:"blog_id" gorm:"not null;index:idx_user_blog,unique;index"`
    Blog      Blog           `json:"-" gorm:"foreignKey:BlogID"`
    CreatedAt time.Time      `json:"created_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}