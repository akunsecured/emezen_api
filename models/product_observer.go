package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductObserver struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	UserID      string             `json:"user_id" bson:"user_id"`
	ProductList []string           `json:"product_list" bson:"product_list"`
}
