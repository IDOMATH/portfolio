package db

import (
	"context"
	"fmt"

	"github.com/IDOMATH/portfolio/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const blogCollection = "blogs"

type BlogStore interface {
	Dropper

	InsertBlogPost(context.Context, *types.BlogPost) (*types.BlogPost, error)
	GetBlogCards(context.Context) ([]*types.BlogCard, error)
	GetBlogById(context.Context, string) (*types.BlogPost, error)
}

type MongoBlogStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewBlogStore(client *mongo.Client, dbName string) *MongoBlogStore {
	return &MongoBlogStore{
		client:     client,
		collection: client.Database(dbName).Collection(blogCollection),
	}
}

func (s *MongoBlogStore) Drop(ctx context.Context) error {
	fmt.Println("--- Dropping blog collection")
	return s.collection.Drop(ctx)
}

func (s *MongoBlogStore) InsertBlogPost(ctx context.Context, blog *types.BlogPost) (*types.BlogPost, error) {
	res, err := s.collection.InsertOne(ctx, blog)
	if err != nil {
		return nil, err
	}
	blog.Id = res.InsertedID.(primitive.ObjectID)
	return blog, nil
}

func (s *MongoBlogStore) GetBlogCards(ctx context.Context) ([]*types.BlogCard, error) {
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var blogCards []*types.BlogCard
	if err := cursor.All(ctx, &blogCards); err != nil {
		return nil, err
	}
	return blogCards, nil
}

func (s *MongoBlogStore) GetBlogById(ctx context.Context, id string) (*types.BlogPost, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var blogPost types.BlogPost
	if err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&blogPost); err != nil {
		return nil, err
	}
	return &blogPost, nil
}
