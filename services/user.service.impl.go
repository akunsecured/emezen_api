package services

import (
	"context"
	"time"

	"github.com/akunsecured/emezen_api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/akunsecured/emezen_api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	userCollection *mongo.Collection
	ctx            context.Context
}

func NewUserService(userCollection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{
		userCollection: userCollection,
		ctx:            ctx,
	}
}

func (u *UserServiceImpl) CreateUser(user *models.User) (*string, error) {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Credits = 0

	result, err := u.userCollection.InsertOne(u.ctx, user)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		var oidHex = oid.Hex()
		return &oidHex, nil
	}
	return nil, utils.ErrInsertedIDIsNotObjectID
}

func (u *UserServiceImpl) GetUser(userId *string) (*models.User, error) {
	var user *models.User
	objID, err := primitive.ObjectIDFromHex(*userId)
	if err != nil {
		return nil, err
	}
	query := bson.D{bson.E{Key: "_id", Value: objID}}
	err = u.userCollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}

func (u *UserServiceImpl) UpdateUser(user *models.User) error {
	filter := bson.D{bson.E{Key: "_id", Value: user.ID}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "first_name", Value: user.FirstName},
		bson.E{Key: "last_name", Value: user.LastName},
		bson.E{Key: "age", Value: user.Age},
		bson.E{Key: "contact_email", Value: user.ContactEmail},
		bson.E{Key: "profile_picture", Value: user.ProfilePicture},
		bson.E{Key: "credits", Value: user.Credits},
		bson.E{Key: "updated_at", Value: time.Now()},
	}}}

	result, _ := u.userCollection.UpdateOne(u.ctx, filter, update)
	if result.MatchedCount != 1 {
		return utils.ErrNotExists
	}
	return nil
}

func (u *UserServiceImpl) DeleteUser(userId *string) error {
	objID, err := primitive.ObjectIDFromHex(*userId)
	if err != nil {
		return err
	}
	filter := bson.D{bson.E{Key: "_id", Value: objID}}
	result, _ := u.userCollection.DeleteOne(u.ctx, filter)
	if result.DeletedCount != 1 {
		return utils.ErrNoMatchedDocumentFoundForDelete
	}
	return nil
}
