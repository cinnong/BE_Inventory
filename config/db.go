package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	// Get MongoDB URI from environment variables
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is required")
	}

	// Get database name from environment variables
	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		dbName = "InventarisKantor" // Default database name
	}

	fmt.Printf("Attempting to connect to MongoDB Atlas...\n")
	fmt.Printf("URI: %s\n", maskPassword(mongoURI))

	// Configure client options with better timeout and retry settings
	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetMaxPoolSize(50).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(30 * time.Second).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(10 * time.Second)

	// Create context with longer timeout for initial connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Failed to create MongoDB client: %v", err)
		log.Fatal("error parsing uri: ", err)
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		log.Fatal("error connecting to MongoDB: ", err)
	}

	DB = client.Database(dbName)
	fmt.Printf("✅ MongoDB Atlas connected successfully! Database: %s\n", dbName)
}

// Helper function to mask password in URI for logging
func maskPassword(uri string) string {
	// Simple masking - replace password with ***
	// This is a basic implementation, you might want to use regex for better masking
	return uri // For now, just return as is since it's for debugging
}
