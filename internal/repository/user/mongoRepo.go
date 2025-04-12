package user

import (
	"chillfix/internal/database"
	"chillfix/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoUserRepository struct {
	db         *database.MongoClient
	collection *mongo.Collection
}

func NewMongoUserRepository(db *database.MongoClient) UserRepository {
	return &mongoUserRepository{
		db:         db,
		collection: db.GetCollection("users"),
	}
}

func (r *mongoUserRepository) Create(ctx context.Context, user *models.User) error {
	if user.Id == "" {
		user.Id = primitive.NewObjectID().Hex()
	}
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *mongoUserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mongoUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := r.collection.FindOne(ctx, bson.M{"email": email})
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *mongoUserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": user.Id},
		bson.M{"$set": user},
	)
	return err
}

func (r *mongoUserRepository) Upsert(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"email": user.Email},
		bson.M{"$set": user},
	)
	if err != nil {
		return err
	}
	fmt.Println("matched count : ", result.MatchedCount)
	if result.MatchedCount == 0 {
		_, err := r.collection.InsertOne(ctx, user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *mongoUserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
