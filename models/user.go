package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName    string             `json:"first_name" bson:"first_name"`
	LastName     string             `json:"last_name" bson:"last_name"`
	Age          int                `json:"age" bson:"age"`
	ContactEmail string             `json:"contact_email" bson:"contact_email"`
	PhoneNumber  int                `json:"phone_number" bson:"phone_number"`
}

type UserCredentials struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	UserID    [12]byte           `json:"user_id" bson:"user_id"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Password  string             `json:"password" bson:"password" validate:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
