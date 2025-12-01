package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Email     string               `json:"email" bson:"email"`
	Name      string               `json:"name" bson:"name"`
	Password  string               `json:"-" bson:"password"`
	Following []Follows            `json:"following,omitempty" bson:"following,omitempty"`
	Followers []Follows            `json:"followers,omitempty" bson:"followers,omitempty"`
	Blogs     []primitive.ObjectID `json:"blogs,omitempty" bson:"blogs,omitempty"`
	Comments  []primitive.ObjectID `json:"comments,omitempty" bson:"comments,omitempty"`
}

// UserResponse represents user data sent to client (without password)
type UserResponse struct {
	ID        primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Email     string               `json:"email" bson:"email"`
	Name      string               `json:"name" bson:"name"`
	Following []Follows            `json:"following,omitempty" bson:"following,omitempty"`
	Followers []Follows            `json:"followers,omitempty" bson:"followers,omitempty"`
	Blogs     []primitive.ObjectID `json:"blogs,omitempty" bson:"blogs,omitempty"`
	Comments  []primitive.ObjectID `json:"comments,omitempty" bson:"comments,omitempty"`
}

// UserCredentials represents login credentials
type UserCredentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserSignupRequest represents signup request
type UserSignupRequest struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}
