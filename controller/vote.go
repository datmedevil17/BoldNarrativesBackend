package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/datmedevil17/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ============ VOTE BUSINESS LOGIC ============

func upvoteBlog(userID, blogID string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}

	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return fmt.Errorf("invalid blog ID: %v", err)
	}

	var existingVote models.Vote
	filter := bson.M{
		"user_id": userObjID,
		"blog_id": blogObjID,
	}
	err = votesCollection.FindOne(context.Background(), filter).Decode(&existingVote)

	if err == nil {
		return errors.New("user has already upvoted this blog")
	}

	if err != mongo.ErrNoDocuments {
		return fmt.Errorf("database error: %v", err)
	}

	vote := models.Vote{
		ID:     primitive.NewObjectID(),
		UserId: userObjID,
		BlogId: blogObjID,
	}

	_, err = votesCollection.InsertOne(context.Background(), vote)
	if err != nil {
		return fmt.Errorf("failed to upvote blog: %v", err)
	}

	return nil
}

func removeUpvote(userID, blogID string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}

	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return fmt.Errorf("invalid blog ID: %v", err)
	}
	filter := bson.M{
		"user_id": userObjID,
		"blog_id": blogObjID,
	}

	result, err := votesCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to remove upvote: %v", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("upvote not found")
	}

	return nil
}

func hasUpvoted(userID, blogID string) (bool, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, fmt.Errorf("invalid user ID: %v", err)
	}

	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return false, fmt.Errorf("invalid blog ID: %v", err)
	}
	filter := bson.M{
		"user_id": userObjID,
		"blog_id": blogObjID,
	}

	var vote models.Vote
	err = votesCollection.FindOne(context.Background(), filter).Decode(&vote)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, fmt.Errorf("database error: %v", err)
	}

	return true, nil
}

func getBlogUpvoteCount(blogID string) (int, error) {
	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return 0, fmt.Errorf("invalid blog ID: %v", err)
	}

	count, err := votesCollection.CountDocuments(context.TODO(), bson.M{"blog_id": blogObjID})
	if err != nil {
		return 0, fmt.Errorf("failed to count upvotes: %v", err)
	}

	return int(count), nil
}

func toggleVote(userID, blogID string) (bool, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, fmt.Errorf("invalid user ID: %v", err)
	}

	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return false, fmt.Errorf("invalid blog ID: %v", err)
	}

	filter := bson.M{
		"user_id": userObjID,
		"blog_id": blogObjID,
	}

	var existingVote models.Vote
	err = votesCollection.FindOne(context.Background(), filter).Decode(&existingVote)

	if err == nil {
		_, err = votesCollection.DeleteOne(context.Background(), filter)
		if err != nil {
			return false, fmt.Errorf("failed to remove vote: %v", err)
		}
		return false, nil
	}

	if err != mongo.ErrNoDocuments {
		return false, fmt.Errorf("database error: %v", err)
	}

	vote := models.Vote{
		ID:     primitive.NewObjectID(),
		UserId: userObjID,
		BlogId: blogObjID,
	}

	_, err = votesCollection.InsertOne(context.Background(), vote)
	if err != nil {
		return false, fmt.Errorf("failed to add vote: %v", err)
	}

	return true, nil
}

// ============ VOTE HTTP HANDLERS ============

func ToggleVoteBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var voteData struct {
		UserID string `json:"user_id"`
		BlogID string `json:"blog_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&voteData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	voted, err := toggleVote(voteData.UserID, voteData.BlogID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := "Vote removed"
	if voted {
		message = "Vote added"
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"voted":   voted,
		"message": message,
	})
}

func CheckVoteStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var voteData struct {
		UserID string `json:"user_id"`
		BlogID string `json:"blog_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&voteData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	voted, err := hasUpvoted(voteData.UserID, voteData.BlogID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	upvoteCount, err := getBlogUpvoteCount(voteData.BlogID)
	if err != nil {
		upvoteCount = 0
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"has_voted":    voted,
		"upvote_count": upvoteCount,
	})
}
