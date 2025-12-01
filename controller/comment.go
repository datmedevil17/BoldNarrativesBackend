package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/datmedevil17/backend/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ============ COMMENT BUSINESS LOGIC ============

func createComment(comment models.Comment) (*models.Comment, error) {
	comment.CreatedAt = time.Now()

	result, err := commentsCollection.InsertOne(context.Background(), comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %v", err)
	}

	comment.ID = result.InsertedID.(primitive.ObjectID)
	fmt.Printf("Comment created successfully with ID: %v\n", comment.ID.Hex())
	return &comment, nil
}

func deleteComment(commentID string) error {
	objectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return fmt.Errorf("invalid comment ID: %v", err)
	}
	filter := bson.M{"_id": objectID}

	result, err := commentsCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %v", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("comment not found")
	}

	return nil
}

func getCommentByID(commentID string) (*models.Comment, error) {
	objectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, fmt.Errorf("invalid comment ID: %v", err)
	}

	var comment models.Comment
	err = commentsCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&comment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("comment not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	return &comment, nil
}

// ============ COMMENT HTTP HANDLERS ============

func CreateComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	result, err := createComment(comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	blogID := params["id"]
	if blogID == "" {
		http.Error(w, "Blog ID is required", http.StatusBadRequest)
		return
	}

	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	cursor, err := commentsCollection.Find(context.TODO(), bson.M{"blog_id": blogObjID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var comments []models.Comment
	if err = cursor.All(context.TODO(), &comments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	commentID := params["id"]
	if commentID == "" {
		http.Error(w, "Comment ID is required", http.StatusBadRequest)
		return
	}

	err := deleteComment(commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Comment deleted successfully"})
}
