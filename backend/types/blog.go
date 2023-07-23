package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// TODO: Maybe make the body html
type BlogPost struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Author      string             `bson:"author" json:"author"`
	PublishedAt time.Time          `bson:"publishedAt" json:"publishedAt"`
	Body        string             `bson:"body" json:"body"`
}

type BlogCard struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Author      string             `bson:"author" json:"author"`
	PublishedAt time.Time          `bson:"publishedAt" json:"publishedAt"`
}
