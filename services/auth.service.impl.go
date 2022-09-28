package services

import (
	"context"
	"github.com/akunsecured/emezen_api/models"
	"github.com/akunsecured/emezen_api/security"
	"github.com/akunsecured/emezen_api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type AuthServiceImpl struct {
	authCollection *mongo.Collection
	ctx            context.Context
}

func NewAuthService(authCollection *mongo.Collection, ctx context.Context) AuthService {
	return &AuthServiceImpl{
		authCollection: authCollection,
		ctx:            ctx,
	}
}

func (a *AuthServiceImpl) CheckIfExistsWithEmail(email string) (*models.UserCredentials, error) {
	var exists *models.UserCredentials
	query := bson.D{bson.E{Key: "email", Value: email}}
	err := a.authCollection.FindOne(a.ctx, query).Decode(&exists)
	return exists, err
}

func (a *AuthServiceImpl) CheckIfExistsWithID(id string) (*models.UserCredentials, error) {
	var exists *models.UserCredentials
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	query := bson.D{bson.E{Key: "_id", Value: objID}}
	err = a.authCollection.FindOne(a.ctx, query).Decode(&exists)
	return exists, err
}

// Register will check if the given email address is already in the database.
// If so, an error will be returned. Otherwise, it will be saved to the database.
func (a *AuthServiceImpl) Register(userCredentials *models.UserCredentials) error {
	exists, err := a.CheckIfExistsWithEmail(userCredentials.Email)
	if exists != nil {
		return utils.ErrEmailIsAlreadyInUse
	}

	if err == mongo.ErrNoDocuments {
		userCredentials.Password, err = security.EncryptPassword(userCredentials.Password)
		if err != nil {
			return err
		}

		userCredentials.CreatedAt = time.Now()
		userCredentials.UpdatedAt = userCredentials.CreatedAt
		userCredentials.ID = primitive.NewObjectID()

		_, err = a.authCollection.InsertOne(a.ctx, userCredentials)
		if err != nil {
			return err
		}

		return nil
	}
	return err
}

// Login will check if the given email is in the database. If not, it will return an error.
// Otherwise, it will check if the password matches with the one in the database. If so, it
// will return a JWT token. Otherwise, it will return an error.
func (a *AuthServiceImpl) Login(userCredentials *models.UserCredentials) (*string, error) {
	exists, err := a.CheckIfExistsWithEmail(userCredentials.Email)
	if err == mongo.ErrNoDocuments {
		return nil, utils.ErrNoAccountWithThisEmail
	}
	if err != nil {
		return nil, err
	}

	err = security.VerifyPassword(exists.Password, userCredentials.Password)
	if err != nil {
		return nil, utils.ErrInvalidPassword
	}

	// TODO: replace with JWT token
	var token = "token"
	return &token, nil
}

// Update will check if the given account is in the database. If not, it will return an error.
// Otherwise, it will update the credentials.
func (a *AuthServiceImpl) Update(userCredentials *models.UserCredentials) error {
	encryptedPassword, err := security.EncryptPassword(userCredentials.Password)
	if err != nil {
		return err
	}

	filter := bson.D{bson.E{Key: "_id", Value: userCredentials.ID}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "email", Value: userCredentials.Email}, bson.E{Key: "password", Value: encryptedPassword}, bson.E{Key: "updated_at", Value: time.Now()}}}}
	result, err := a.authCollection.UpdateOne(a.ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount != 1 {
		return utils.ErrNotExists
	}
	return nil
}
