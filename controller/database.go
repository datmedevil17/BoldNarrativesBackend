package controller

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client             *mongo.Client
	usersCollection    *mongo.Collection
	blogsCollection    *mongo.Collection
	commentsCollection *mongo.Collection
	votesCollection    *mongo.Collection
	followsCollection  *mongo.Collection
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Get MongoDB configuration from environment variables
	connectionString := os.Getenv("MONGO_URI")
	if connectionString == "" {
		log.Fatal("MONGO_URI environment variable is required")
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		log.Fatal("MONGO_DB_NAME environment variable is required")
	}

	// Connect to MongoDB
	clientOption := options.Client().ApplyURI(connectionString)
	var err error
	client, err = mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	database := client.Database(dbName)
	usersCollection = database.Collection("users")
	blogsCollection = database.Collection("blogs")
	commentsCollection = database.Collection("comments")
	votesCollection = database.Collection("votes")
	followsCollection = database.Collection("follows")

	log.Println("Connected to MongoDB successfully!")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetUsersCollection() *mongo.Collection {
	return usersCollection
}

func GetBlogsCollection() *mongo.Collection {
	return blogsCollection
}

func GetCommentsCollection() *mongo.Collection {
	return commentsCollection
}

func GetVotesCollection() *mongo.Collection {
	return votesCollection
}

func GetFollowsCollection() *mongo.Collection {
	return followsCollection
}

func DisconnectDB() error {
	if client != nil {
		return client.Disconnect(context.TODO())
	}
	return nil
}
