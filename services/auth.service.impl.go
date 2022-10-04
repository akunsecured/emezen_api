package services

import (
	"context"
	"github.com/akunsecured/emezen_api/models"
	"github.com/akunsecured/emezen_api/security"
	"github.com/akunsecured/emezen_api/utils"
	"github.com/form3tech-oss/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type AuthServiceImpl struct {
	authCollection *mongo.Collection
	userService    UserService
	ctx            context.Context
}

func NewAuthService(authCollection *mongo.Collection, userService UserService, ctx context.Context) AuthService {
	return &AuthServiceImpl{
		authCollection: authCollection,
		userService:    userService,
		ctx:            ctx,
	}
}

// CheckIfExistsWithEmail will check if there are any credentials with the given email address.
// If so, it will be returned. Otherwise, it will return an error.
func (a *AuthServiceImpl) CheckIfExistsWithEmail(email string) (*models.UserCredentials, error) {
	var exists *models.UserCredentials
	query := bson.D{bson.E{Key: "email", Value: email}}
	err := a.authCollection.FindOne(a.ctx, query).Decode(&exists)
	return exists, err
}

// CheckIfExistsWithID will check if there are any user credentials with the given ID.
// If so, it will be returned. Otherwise, it will return an error.
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
func (a *AuthServiceImpl) Register(userDataWithCredentials *models.UserDataWithCredentials) (*models.WrappedToken, error) {
	var userCredentials = userDataWithCredentials.Credentials

	exists, err := a.CheckIfExistsWithEmail(userCredentials.Email)
	if exists != nil {
		return nil, utils.ErrEmailIsAlreadyInUse
	}

	if err == mongo.ErrNoDocuments {
		var userData = userDataWithCredentials.UserData
		if len(userData.ContactEmail) == 0 {
			userData.ContactEmail = userCredentials.Email
		}
		userId, err := a.userService.CreateUser(&userData)
		if err != nil {
			return nil, err
		}

		userCredentials.Password, err = security.EncryptPassword(userCredentials.Password)
		if err != nil {
			return nil, err
		}

		userCredentials.CreatedAt = time.Now()
		userCredentials.UpdatedAt = userCredentials.CreatedAt
		userCredentials.ID = primitive.NewObjectID()
		userCredentials.UserID = *userId

		_, err = a.authCollection.InsertOne(a.ctx, userCredentials)
		if err != nil {
			return nil, err
		}

		userData.ID, err = primitive.ObjectIDFromHex(*userId)
		if err != nil {
			return nil, err
		}

		wrappedToken, err := security.CreateAccessAndRefreshTokens(userData)
		if err != nil {
			return nil, err
		}

		return wrappedToken, nil
	}
	return nil, err
}

// Login will check if the given email is in the database. If not, it will return an error.
// Otherwise, it will check if the password matches with the one in the database. If so, it
// will return a JWT token. Otherwise, it will return an error.
func (a *AuthServiceImpl) Login(userCredentials *models.UserCredentials) (*models.WrappedToken, error) {
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

	user, err := a.userService.GetUser(&exists.UserID)
	if err != nil {
		return nil, mongo.ErrNoDocuments
	}

	wrappedToken, err := security.CreateAccessAndRefreshTokens(*user)
	if err != nil {
		return nil, err
	}

	return wrappedToken, nil
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

// NewAccessToken will try to parse the incoming token string and create a new access token from
// the user ID given in the refresh token's claims
func (a *AuthServiceImpl) NewAccessToken(claims *jwt.MapClaims) (*string, error) {
	user, err := a.CurrentUser(claims)
	if err != nil {
		return nil, err
	}

	newToken, err := security.NewAccessToken(*user)
	if err != nil {
		return nil, err
	}

	return &newToken, nil
}

func (a *AuthServiceImpl) CurrentUser(claims *jwt.MapClaims) (*models.User, error) {
	userId := (*claims)["sub"].(string)

	user, err := a.userService.GetUser(&userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
