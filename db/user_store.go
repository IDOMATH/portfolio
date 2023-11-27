package db

import (
	"context"
	"fmt"
	"github.com/IDOMATH/portfolio/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const userCollection = "users"

type UserStore interface {
	Dropper

	InsertUser(context.Context, *types.User) (*types.User, error)
}

type MongoUserStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewUserStore(client *mongo.Client, dbName string) *MongoBlogStore {
	return &MongoBlogStore{
		client:     client,
		collection: client.Database(dbName).Collection(userCollection),
	}
}

func HashPassword(password string) (string, error) {
	//TODO: do more research into some various cost recommendations for production
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return "", err
	}

	return string(passwordHash), nil
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("--- Dropping blog collection")
	return s.collection.Drop(ctx)
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	var err error
	user.Password, err = HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	res, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Id = res.InsertedID.(primitive.ObjectID)
	return user, nil
}
