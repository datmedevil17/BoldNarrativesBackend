package models

import (
    "time"
    "gorm.io/gorm"
)

type Follows struct {
    FollowerID  uint           `json:"follower_id" gorm:"primaryKey;index"`
    FollowingID uint           `json:"following_id" gorm:"primaryKey;index"`
    Follower    User           `json:"follower" gorm:"foreignKey:FollowerID"`
    Following   User           `json:"following" gorm:"foreignKey:FollowingID"`
    CreatedAt   time.Time      `json:"created_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type FollowResponse struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}