package services

import (
	"context"
	"github.com/akunsecured/emezen_api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

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
	// TODO: Fix the update of the profile
	/*
		filter := bson.D{bson.E{Key: "name", Value: user.Name}}
		update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "name", Value: user.Name}, bson.E{Key: "email", Value: user.Email}, bson.E{Key: "age", Value: user.Age}}}}
		result, _ := u.userCollection.UpdateOne(u.ctx, filter, update)
		if result.MatchedCount != 1 {
			return errors.New("No matched document found for update")
		}
	*/
	return utils.ErrUnimplementedMethod
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
