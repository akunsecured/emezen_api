package services

import (
	"context"
	"errors"
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

func (u *UserServiceImpl) GetUser(name *string) (*models.User, error) {
	var user *models.User
	// db.collection.find({name: name})
	query := bson.D{bson.E{Key: "name", Value: name}}
	err := u.userCollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}

func (u *UserServiceImpl) GetAll() ([]*models.User, error) {
	var users []*models.User
	cursor, err := u.userCollection.Find(u.ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(u.ctx) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	err = cursor.Close(u.ctx)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, errors.New("Documents not found")
	}
	return users, nil
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
	return nil
}

func (u *UserServiceImpl) DeleteUser(name *string) error {
	filter := bson.D{bson.E{Key: "name", Value: name}}
	result, _ := u.userCollection.DeleteOne(u.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("No matched document found for update")
	}
	return nil
}
