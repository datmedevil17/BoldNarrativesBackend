package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Follows struct {
	FollowerId  primitive.ObjectID `json:"follower_id" bson:"follower_id"`
	FollowingId primitive.ObjectID `json:"following_id" bson:"following_id"`
}

// FollowRequest represents a follow/unfollow request
type FollowRequest struct {
	FollowerId  primitive.ObjectID `json:"follower_id" binding:"required"`
	FollowingId primitive.ObjectID `json:"following_id" binding:"required"`
}

// FollowResponse represents follow status response
type FollowResponse struct {
	IsFollowing bool `json:"is_following"`
}

// FollowersResponse represents a list of followers with user info
type FollowersResponse struct {
	Followers []UserResponse `json:"followers"`
	Count     int            `json:"count"`
}

// FollowingResponse represents a list of following with user info
type FollowingResponse struct {
	Following []UserResponse `json:"following"`
	Count     int            `json:"count"`
}
