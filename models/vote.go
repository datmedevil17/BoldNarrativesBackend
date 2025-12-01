package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vote struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
	BlogId primitive.ObjectID `json:"blog_id" bson:"blog_id"`
}

// VoteRequest represents a vote/unvote request
type VoteRequest struct {
	UserId primitive.ObjectID `json:"user_id" binding:"required"`
	BlogId primitive.ObjectID `json:"blog_id" binding:"required"`
}

// VoteResponse represents vote status response
type VoteResponse struct {
	HasVoted    bool `json:"has_voted"`
	UpvoteCount int  `json:"upvote_count"`
}
