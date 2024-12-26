package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kihyun1998/prisma-market/prisma-user-service/internal/models"
)

type UserRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewUserRepository(mongoURI string) (*UserRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database("prisma_market")
	collection := db.Collection("users")

	return &UserRepository{
		db:         db,
		collection: collection,
	}, nil
}

func (r *UserRepository) CreateProfile(ctx context.Context, profile *models.UserProfile) error {
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()
	profile.Status = "active"

	result, err := r.collection.InsertOne(ctx, profile)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("profile already exists")
		}
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		profile.ID = oid
	}

	return nil
}

func (r *UserRepository) GetProfileByID(ctx context.Context, id primitive.ObjectID) (*models.UserProfile, error) {
	var profile models.UserProfile
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&profile)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *UserRepository) GetProfileByAuthID(ctx context.Context, authID primitive.ObjectID) (*models.UserProfile, error) {
	var profile models.UserProfile
	err := r.collection.FindOne(ctx, bson.M{"auth_id": authID}).Decode(&profile)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *UserRepository) GetProfileByUsername(ctx context.Context, username string) (*models.UserProfile, error) {
	var profile models.UserProfile
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&profile)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now()
	result, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("profile not found")
	}
	return nil
}

func (r *UserRepository) DeleteProfile(ctx context.Context, id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"status":     "inactive",
			"updated_at": time.Now(),
		},
	}
	result, err := r.collection.UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("profile not found")
	}
	return nil
}

func (r *UserRepository) SearchProfiles(ctx context.Context, query string, limit int64) ([]*models.UserProfile, error) {
	filter := bson.M{
		"$text": bson.M{
			"$search": query,
		},
		"status": "active",
	}

	findOptions := options.Find().
		SetLimit(limit).
		SetSort(bson.D{{Key: "score", Value: bson.M{"$meta": "textScore"}}})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var profiles []*models.UserProfile
	if err = cursor.All(ctx, &profiles); err != nil {
		return nil, err
	}

	return profiles, nil
}

func (r *UserRepository) Close(ctx context.Context) error {
	return r.db.Client().Disconnect(ctx)
}
