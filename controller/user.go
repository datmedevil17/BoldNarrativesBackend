package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/datmedevil17/backend/models"
	"github.com/datmedevil17/backend/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ============ USER BUSINESS LOGIC ============

func userSignUp(user models.User) (*models.User, error) {
	var existingUser models.User
	filter := bson.M{"email": user.Email}
	err := usersCollection.FindOne(context.TODO(), filter).Decode(&existingUser)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}
	if err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("database error: %v", err)
	}
	result, err := usersCollection.InsertOne(context.TODO(), user)
	check(err)
	user.ID = result.InsertedID.(primitive.ObjectID)
	user.Password = ""
	fmt.Printf("User signed up successfully with ID: %v\n", user.ID.Hex())
	return &user, nil
}

func userSignIn(email, password string) (*models.User, error) {
	var user models.User
	err := usersCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid email or password")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}
	if user.Password != password {
		return nil, errors.New("invalid email or password")
	}
	user.Password = ""
	fmt.Printf("User signed in successfully: %v\n", user.Email)
	return &user, nil
}

func getUserByID(userID string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	check(err)
	var user models.User
	filter := bson.M{"_id": objectID}
	err = usersCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}
	user.Password = ""
	return &user, nil
}

// ============ USER HTTP HANDLERS ============

func UserSignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	result, err := userSignUp(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func UserSignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	user, err := userSignIn(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.Hex(), user.Email)
	if err != nil {
		http.Error(w, "Failed to generate authentication token", http.StatusInternalServerError)
		return
	}

	// Return user data and token
	response := map[string]interface{}{
		"user":  user,
		"token": token,
	}
	json.NewEncoder(w).Encode(response)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	userID := params["id"]
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	result, err := getUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func ViewUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	userID := params["id"]
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	result, err := getUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func GetAuthenticatedUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"user_id": userID})
}

func GetUserProfileByQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "User ID query parameter is required", http.StatusBadRequest)
		return
	}

	result, err := getUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func ProtectedRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	userContext, err := utils.GetUserFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Access granted to protected resource",
		"user_id": userContext.UserID,
		"email":   userContext.Email,
	})
}
