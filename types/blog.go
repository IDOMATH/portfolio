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
	ImageName   string             `bson:"imageName" json:"imageName"`
}

// TODO: Maybe we don't need this.  Could just be a function to get these values
type BlogCard struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Author      string             `bson:"author" json:"author"`
	PublishedAt time.Time          `bson:"publishedAt" json:"publishedAt"`
	ImageName   string             `bson:"imageName" json:"imageName"`
	secretField string
}
