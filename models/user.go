package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	LastVisitAt time.Time          `json:"lastVisitAt" bson:"lastVisitAt"`
	Session     Session            `json:"session" bson:"session,omitempty"`
}
