package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID        primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string               `json:"title" bson:"title"`
	AuthorId  primitive.ObjectID   `json:"author_id" bson:"author_id"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
	Content   string               `json:"content" bson:"content"`
	Genre     string               `json:"genre" bson:"genre"`
	Views     int                  `json:"views" bson:"views"`
	Votes     []Vote               `json:"votes,omitempty" bson:"votes,omitempty"`
	Comments  []primitive.ObjectID `json:"comments,omitempty" bson:"comments,omitempty"`
}

// BlogRequest represents a blog creation/update request
type BlogRequest struct {
	Title    string             `json:"title" binding:"required"`
	AuthorId primitive.ObjectID `json:"author_id" binding:"required"`
	Content  string             `json:"content" binding:"required"`
	Genre    string             `json:"genre" binding:"required"`
}

// BlogResponse represents a blog response with additional computed fields
type BlogResponse struct {
	ID         primitive.ObjectID   `json:"_id,omitempty"`
	Title      string               `json:"title"`
	AuthorId   primitive.ObjectID   `json:"author_id"`
	AuthorName string               `json:"author_name,omitempty"`
	CreatedAt  time.Time            `json:"created_at"`
	Content    string               `json:"content"`
	Genre      string               `json:"genre"`
	Views      int                  `json:"views"`
	VoteCount  int                  `json:"vote_count"`
	Comments   []primitive.ObjectID `json:"comments,omitempty"`
}

// BlogFilter represents filtering options for blog queries
type BlogFilter struct {
	Genre    string `json:"genre,omitempty"`
	Search   string `json:"search,omitempty"`
	AuthorId string `json:"author_id,omitempty"`
	SortBy   string `json:"sort_by,omitempty"` // "time", "views", "trending"
	Asc      bool   `json:"asc,omitempty"`
	Skip     int    `json:"skip,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}
