package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"unique;not null;index"`
	Name      string         `json:"name" gorm:"not null"`
	Password  string         `json:"-" gorm:"not null"`
	Blogs     []Blog         `json:"blogs,omitempty" gorm:"foreignKey:AuthorID"`
	Comments  []Comment      `json:"comments,omitempty" gorm:"foreignKey:AuthorID"`
	Following []Follows      `json:"following,omitempty" gorm:"foreignKey:FollowerID"`
	Followers []Follows      `json:"followers,omitempty" gorm:"foreignKey:FollowingID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserResponse struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Email string `json:"email" gorm:"unique;not null;index"`
	Name  string `json:"name" gorm:"not null"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
	}
}
