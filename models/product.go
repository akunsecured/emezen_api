package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	SellerID  string             `json:"seller_id" bson:"seller_id" validate:"required"`
	Name      string             `json:"name" bson:"name" validate:"required,min=1,max=50"`
	Price     float32            `json:"price" bson:"price" validate:"required,min=0.1"`
	Images    []string           `json:"images" bson:"images"`
	Details   string             `json:"details" bson:"details" validate:"required,min=1,max=500"`
	Quantity  int32              `json:"quantity" bson:"quantity" validate:"required,min=1,max=100"`
	Category  Category           `json:"category" bson:"category"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
}
