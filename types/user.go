package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// TODO: probably hide the password or something like that.  Need to double check on this.
type User struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"password"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
