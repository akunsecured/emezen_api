package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName    string             `json:"first_name" bson:"first_name" validate:"required,min=1,max=50"`
	LastName     string             `json:"last_name" bson:"last_name" validate:"required,min=1,max=50"`
	Age          int                `json:"age" bson:"age" validate:"required,min=13,max=100"`
	ContactEmail string             `json:"contact_email" bson:"contact_email"`
	PhoneNumber  int                `json:"phone_number" bson:"phone_number"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserCredentials struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Password  string             `json:"password" bson:"password" validate:"required,min=8,max=64"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserDataWithCredentials struct {
	UserData    User            `json:"user_data"`
	Credentials UserCredentials `json:"credentials"`
}
