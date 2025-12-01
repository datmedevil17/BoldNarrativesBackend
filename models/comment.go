package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AuthorId  primitive.ObjectID `json:"author_id" bson:"author_id"`
	Comment   string             `json:"comment" bson:"comment"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	BlogId    primitive.ObjectID `json:"blog_id" bson:"blog_id"`
}

// CommentRequest represents a comment creation request
type CommentRequest struct {
	AuthorId primitive.ObjectID `json:"author_id" binding:"required"`
	Comment  string             `json:"comment" binding:"required"`
	BlogId   primitive.ObjectID `json:"blog_id" binding:"required"`
}

// CommentResponse represents a comment with author information
type CommentResponse struct {
	ID         primitive.ObjectID `json:"_id,omitempty"`
	AuthorId   primitive.ObjectID `json:"author_id"`
	AuthorName string             `json:"author_name,omitempty"`
	Comment    string             `json:"comment"`
	CreatedAt  time.Time          `json:"created_at"`
	BlogId     primitive.ObjectID `json:"blog_id"`
}
