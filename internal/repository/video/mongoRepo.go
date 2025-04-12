package video

import (
	"chillfix/internal/database"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoVideoRepository struct {
	db         *database.MongoClient
	collection *mongo.Collection
}
