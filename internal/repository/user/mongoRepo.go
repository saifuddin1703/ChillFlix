package user

import (
	"chillfix/internal/database"
	"chillfix/models"
	"context"
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

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	userId := result.InsertedID.(primitive.ObjectID)
	user.SetID(userId)
	return nil
}

func (r *mongoUserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mongoUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mongoUserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	return err
}

func (r *mongoUserRepository) Delete(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectId})
	return err
}
