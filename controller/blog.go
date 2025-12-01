package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/datmedevil17/backend/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ============ BLOG BUSINESS LOGIC ============

func createBlog(blog models.Blog) (*models.Blog, error) {
	blog.CreatedAt = time.Now()
	blog.Views = 0

	result, err := blogsCollection.InsertOne(context.Background(), blog)
	if err != nil {
		return nil, fmt.Errorf("failed to create blog: %v", err)
	}

	blog.ID = result.InsertedID.(primitive.ObjectID)
	fmt.Printf("Blog created successfully with ID: %v\n", blog.ID.Hex())
	return &blog, nil
}

func deleteBlog(blogID string) error {
	return deleteBlogCascade(blogID)
}

func getBlogByID(blogID string) (*models.Blog, error) {
	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, fmt.Errorf("invalid blog ID: %v", err)
	}

	var blog models.Blog
	filter := bson.M{"_id": objectID}
	err = blogsCollection.FindOne(context.TODO(), filter).Decode(&blog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("blog not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	return &blog, nil
}

func listBlogsByTime(asc bool, skip, limit int) ([]models.Blog, error) {
	sortOrder := 1
	if !asc {
		sortOrder = -1
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: sortOrder}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := blogsCollection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find blogs: %v", err)
	}
	defer cursor.Close(context.TODO())

	var blogs []models.Blog
	if err = cursor.All(context.TODO(), &blogs); err != nil {
		return nil, fmt.Errorf("failed to decode blogs: %v", err)
	}

	return blogs, nil
}

func listBlogsByViews(skip, limit int) ([]models.Blog, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "views", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := blogsCollection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find blogs: %v", err)
	}
	defer cursor.Close(context.TODO())

	var blogs []models.Blog
	if err = cursor.All(context.TODO(), &blogs); err != nil {
		return nil, fmt.Errorf("failed to decode blogs: %v", err)
	}

	return blogs, nil
}

func getTrendingBlogs(limit int) ([]models.Blog, error) {
	return getTrendingBlogsAdvanced(limit)
}

func getBlogCount(genre, search, author string) (int, error) {
	filter := bson.M{}

	if genre != "" {
		filter["genre"] = genre
	}

	if search != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": search, "$options": "i"}},
			{"content": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	if author != "" {
		authorObjID, err := primitive.ObjectIDFromHex(author)
		if err != nil {
			return 0, fmt.Errorf("invalid author ID: %v", err)
		}
		filter["author_id"] = authorObjID
	}

	count, err := blogsCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count blogs: %v", err)
	}

	return int(count), nil
}

func listBlogsFiltered(genre, search, author string, sortBy string, asc bool, skip, limit int) ([]models.Blog, error) {
	filter := bson.M{}

	if genre != "" {
		filter["genre"] = genre
	}

	if search != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": search, "$options": "i"}},
			{"content": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	if author != "" {
		authorObjID, err := primitive.ObjectIDFromHex(author)
		if err != nil {
			return nil, fmt.Errorf("invalid author ID: %v", err)
		}
		filter["author_id"] = authorObjID
	}

	sortOrder := 1
	if !asc {
		sortOrder = -1
	}

	sortField := "created_at"
	if sortBy == "views" {
		sortField = "views"
	} else if sortBy == "trending" {
		sortField = "views"
		sortOrder = -1
	}

	opts := options.Find().
		SetSort(bson.D{{Key: sortField, Value: sortOrder}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := blogsCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find blogs: %v", err)
	}
	defer cursor.Close(context.TODO())

	var blogs []models.Blog
	if err = cursor.All(context.TODO(), &blogs); err != nil {
		return nil, fmt.Errorf("failed to decode blogs: %v", err)
	}

	return blogs, nil
}

func calculateTrendingScore(blog models.Blog, upvoteCount int) float64 {
	now := time.Now()
	ageInHours := now.Sub(blog.CreatedAt).Hours()
	if ageInHours < 1 {
		ageInHours = 1
	}

	score := float64(blog.Views+upvoteCount) / ageInHours
	return score
}

func getTrendingBlogsAdvanced(limit int) ([]models.Blog, error) {
	cursor, err := blogsCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find blogs: %v", err)
	}
	defer cursor.Close(context.TODO())

	var blogs []models.Blog
	if err = cursor.All(context.TODO(), &blogs); err != nil {
		return nil, fmt.Errorf("failed to decode blogs: %v", err)
	}

	type blogWithScore struct {
		Blog  models.Blog
		Score float64
	}

	var blogsWithScores []blogWithScore
	for _, blog := range blogs {
		upvoteCount, err := getBlogUpvoteCount(blog.ID.Hex())
		if err != nil {
			upvoteCount = 0
		}

		score := calculateTrendingScore(blog, upvoteCount)
		blogsWithScores = append(blogsWithScores, blogWithScore{
			Blog:  blog,
			Score: score,
		})
	}

	for i := 0; i < len(blogsWithScores); i++ {
		for j := i + 1; j < len(blogsWithScores); j++ {
			if blogsWithScores[j].Score > blogsWithScores[i].Score {
				blogsWithScores[i], blogsWithScores[j] = blogsWithScores[j], blogsWithScores[i]
			}
		}
	}

	var result []models.Blog
	for i, item := range blogsWithScores {
		if i >= limit {
			break
		}
		result = append(result, item.Blog)
	}

	return result, nil
}

func updateBlog(blogID string, updates bson.M) (*models.Blog, error) {
	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, fmt.Errorf("invalid blog ID: %v", err)
	}

	delete(updates, "_id")
	delete(updates, "created_at")
	delete(updates, "views")

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": updates}

	result := blogsCollection.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, errors.New("blog not found")
		}
		return nil, fmt.Errorf("failed to update blog: %v", result.Err())
	}

	var updatedBlog models.Blog
	if err := result.Decode(&updatedBlog); err != nil {
		return nil, fmt.Errorf("failed to decode updated blog: %v", err)
	}

	return &updatedBlog, nil
}

func deleteBlogCascade(blogID string) error {
	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return fmt.Errorf("invalid blog ID: %v", err)
	}

	_, err = votesCollection.DeleteMany(context.TODO(), bson.M{"blog_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete blog votes: %v", err)
	}

	_, err = commentsCollection.DeleteMany(context.TODO(), bson.M{"blog_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete blog comments: %v", err)
	}

	result, err := blogsCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete blog: %v", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("blog not found")
	}

	return nil
}

func incrementBlogViews(blogID string) {
	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		log.Printf("incrementBlogViews: invalid blog ID %s: %v", blogID, err)
		return
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$inc": bson.M{"views": 1}}

	_, err = blogsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("incrementBlogViews: failed to increment views for blog %s: %v", blogID, err)
		return
	}
}

// ============ BLOG HTTP HANDLERS ============

func CreateBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var blog models.Blog
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	result, err := createBlog(blog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func DeleteBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	blogID := params["id"]
	if blogID == "" {
		http.Error(w, "Blog ID is required", http.StatusBadRequest)
		return
	}

	err := deleteBlog(blogID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Blog deleted successfully")
}

func GetBlogByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	blogID := params["id"]
	if blogID == "" {
		http.Error(w, "Blog ID is required", http.StatusBadRequest)
		return
	}

	result, err := getBlogByID(blogID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	incrementBlogViews(blogID)
	json.NewEncoder(w).Encode(result)
}

func ListBlogsByTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	ascParam := r.URL.Query().Get("asc")
	skipParam := r.URL.Query().Get("skip")
	limitParam := r.URL.Query().Get("limit")

	asc := ascParam == "true"
	skip, _ := strconv.Atoi(skipParam)
	limit, _ := strconv.Atoi(limitParam)

	if limit == 0 {
		limit = 10
	}

	result, err := listBlogsByTime(asc, skip, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func ListBlogsByViews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	skipParam := r.URL.Query().Get("skip")
	limitParam := r.URL.Query().Get("limit")

	skip, _ := strconv.Atoi(skipParam)
	limit, _ := strconv.Atoi(limitParam)

	if limit == 0 {
		limit = 10
	}

	result, err := listBlogsByViews(skip, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func GetTrendingBlogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	limitParam := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitParam)

	if limit == 0 {
		limit = 10
	}

	result, err := getTrendingBlogs(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func GetBlogCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	genre := r.URL.Query().Get("genre")
	search := r.URL.Query().Get("search")
	author := r.URL.Query().Get("author")

	count, err := getBlogCount(genre, search, author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"total": count})
}

func ListBlogsAdvanced(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	genre := r.URL.Query().Get("genre")
	search := r.URL.Query().Get("search")
	author := r.URL.Query().Get("author")
	sortBy := r.URL.Query().Get("sort_by")
	ascParam := r.URL.Query().Get("asc")
	skipParam := r.URL.Query().Get("skip")
	limitParam := r.URL.Query().Get("limit")

	asc := ascParam == "true"
	skip, _ := strconv.Atoi(skipParam)
	limit, _ := strconv.Atoi(limitParam)

	if limit == 0 {
		limit = 10
	}

	if sortBy == "" {
		sortBy = "time"
	}

	result, err := listBlogsFiltered(genre, search, author, sortBy, asc, skip, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func GetAdvancedTrendingBlogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	limitParam := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitParam)

	if limit == 0 {
		limit = 10
	}

	result, err := getTrendingBlogsAdvanced(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func UpdateBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	blogID := params["id"]
	if blogID == "" {
		http.Error(w, "Blog ID is required", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	bsonUpdates := bson.M{}
	for key, value := range updates {
		bsonUpdates[key] = value
	}

	result, err := updateBlog(blogID, bsonUpdates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func GetTotalBlogCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var filters struct {
		Genre  string `json:"genre"`
		Search string `json:"search"`
		Author string `json:"author"`
	}

	if err := json.NewDecoder(r.Body).Decode(&filters); err != nil {
		filters = struct {
			Genre  string `json:"genre"`
			Search string `json:"search"`
			Author string `json:"author"`
		}{}
	}

	count, err := getBlogCount(filters.Genre, filters.Search, filters.Author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"total": count})
}

func SortBlogsByTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	params := mux.Vars(r)
	userID := params["id"]

	var sortData struct {
		Asc   bool `json:"asc"`
		Skip  int  `json:"skip"`
		Limit int  `json:"limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&sortData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if sortData.Limit == 0 {
		sortData.Limit = 10
	}

	var result []models.Blog
	var err error

	if userID != "" {
		result, err = listBlogsFiltered("", "", userID, "time", sortData.Asc, sortData.Skip, sortData.Limit)
	} else {
		result, err = listBlogsByTime(sortData.Asc, sortData.Skip, sortData.Limit)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func SortBlogsByViews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var sortData struct {
		Skip  int `json:"skip"`
		Limit int `json:"limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&sortData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if sortData.Limit == 0 {
		sortData.Limit = 10
	}

	result, err := listBlogsByViews(sortData.Skip, sortData.Limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func IncrementViewCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")

	var viewData struct {
		BlogID string `json:"blog_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&viewData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	incrementBlogViews(viewData.BlogID)

	json.NewEncoder(w).Encode(map[string]string{"message": "View count incremented"})
}
