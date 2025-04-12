package database

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var (
	instance *MongoClient
	once     sync.Once
	mu       sync.Mutex
)

// GetMongoClient returns the singleton instance of MongoClient
func GetMongoClient(uri string, dbName string) (*MongoClient, error) {
	mu.Lock()
	defer mu.Unlock()

	var err error
	once.Do(func() {
		instance, err = newMongoClient(uri, dbName)
	})

	return instance, err
}

// newMongoClient creates a new MongoDB client (private constructor)
func newMongoClient(uri string, dbName string) (*MongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	db := client.Database(dbName)
	return &MongoClient{
		Client: client,
		DB:     db,
	}, nil
}

func (m *MongoClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

// GetCollection returns a mongo collection
func (m *MongoClient) GetCollection(name string) *mongo.Collection {
	return m.DB.Collection(name)
}
