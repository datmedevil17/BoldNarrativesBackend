package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/datmedevil17/backend/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ============ FOLLOW BUSINESS LOGIC ============

func follow(followerID, followingID string) error {
	followerObjID, err := primitive.ObjectIDFromHex(followerID)
	if err != nil {
		return fmt.Errorf("invalid followerID: %v", err)
	}
	followingObjID, err := primitive.ObjectIDFromHex(followingID)
	if err != nil {
		return fmt.Errorf("invalid followingID: %v", err)
	}
	var follow models.Follows
	follow.FollowerId = followerObjID
	follow.FollowingId = followingObjID
	_, err = followsCollection.InsertOne(context.Background(), follow)
	if err != nil {
		return fmt.Errorf("failed to insert follow: %v", err)
	}
	return nil
}

func followCheck(followerID, followingID string) error {
	followerObjID, err := primitive.ObjectIDFromHex(followerID)
	if err != nil {
		return fmt.Errorf("invalid followerID: %v", err)
	}
	followingObjID, err := primitive.ObjectIDFromHex(followingID)
	if err != nil {
		return fmt.Errorf("invalid followingID: %v", err)
	}

	var follow models.Follows
	filter := bson.M{
		"follower_id":  followerObjID,
		"following_id": followingObjID,
	}
	err = followsCollection.FindOne(context.Background(), filter).Decode(&follow)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("not following")
		}
		return fmt.Errorf("failed to check follow: %v", err)
	}
	return nil
}

func unfollow(followerID, followingID string) error {
	followerObjID, err := primitive.ObjectIDFromHex(followerID)
	if err != nil {
		return fmt.Errorf("invalid followerID: %v", err)
	}
	followingObjID, err := primitive.ObjectIDFromHex(followingID)
	if err != nil {
		return fmt.Errorf("invalid followingID: %v", err)
	}
	filter := bson.M{
		"follower_id":  followerObjID,
		"following_id": followingObjID,
	}

	_, err = followsCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to unfollow: %v", err)
	}
	return nil
}

func showFollowers(userID string) ([]models.Follows, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid userID: %v", err)
	}
	filter := bson.M{
		"following_id": userObjID,
	}

	cursor, err := followsCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find followers: %v", err)
	}
	defer cursor.Close(context.Background())

	var followers []models.Follows
	for cursor.Next(context.Background()) {
		var follow models.Follows
		if err := cursor.Decode(&follow); err != nil {
			return nil, fmt.Errorf("failed to decode follower: %v", err)
		}
		followers = append(followers, follow)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return followers, nil
}

func showFollowing(userID string) ([]models.Follows, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid userID: %v", err)
	}
	filter := bson.M{
		"follower_id": userObjID,
	}

	cursor, err := followsCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find following: %v", err)
	}
	defer cursor.Close(context.Background())

	var following []models.Follows
	for cursor.Next(context.Background()) {
		var follow models.Follows
		if err := cursor.Decode(&follow); err != nil {
			return nil, fmt.Errorf("failed to decode following: %v", err)
		}
		following = append(following, follow)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return following, nil
}

// ============ FOLLOW HTTP HANDLERS ============

func Follow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var followData struct {
		FollowerID  string `json:"follower_id"`
		FollowingID string `json:"following_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&followData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	err := follow(followData.FollowerID, followData.FollowingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Successfully followed")
}

func Unfollow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	var followData struct {
		FollowerID  string `json:"follower_id"`
		FollowingID string `json:"following_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&followData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	err := unfollow(followData.FollowerID, followData.FollowingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Successfully unfollowed")
}

func CheckFollow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	followerID := params["follower_id"]
	followingID := params["following_id"]

	if followerID == "" || followingID == "" {
		http.Error(w, "Both follower_id and following_id are required", http.StatusBadRequest)
		return
	}

	err := followCheck(followerID, followingID)
	isFollowing := err == nil

	json.NewEncoder(w).Encode(map[string]bool{"is_following": isFollowing})
}

func ShowFollowers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	userID := params["user_id"]
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	result, err := showFollowers(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func ShowFollowing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	userID := params["user_id"]
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	result, err := showFollowing(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}
