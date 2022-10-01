package main

import (
	"context"
	"fmt"
	"github.com/akunsecured/emezen_api/controllers"
	"github.com/akunsecured/emezen_api/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"time"
)

var (
	server         *gin.Engine
	ctx            context.Context
	mongoClient    *mongo.Client
	mongoDatabase  *mongo.Database
	userCollection *mongo.Collection
	userService    services.UserService
	userController controllers.UserController
	authCollection *mongo.Collection
	authService    services.AuthService
	authController controllers.AuthController
	err            error
)

// This function runs before the main()
func init() {
	fmt.Println("Connecting to MongoDB...")

	ctx = context.TODO()

	mongoConnection := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoClient, err = mongo.Connect(ctx, mongoConnection)
	if err != nil {
		log.Fatal(err)
	}
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	mongoDatabase = mongoClient.Database("emezendb")

	userCollection = mongoDatabase.Collection("users")
	userService = services.NewUserService(userCollection, ctx)
	userController = controllers.NewUserController(userService)

	authCollection = mongoDatabase.Collection("credentials")
	authService = services.NewAuthService(authCollection, userService, ctx)
	authController = controllers.NewAuthController(authService)

	server = gin.Default()
}

func main() {
	defer func(mongoClient *mongo.Client, ctx context.Context) {
		err := mongoClient.Disconnect(ctx)
		if err != nil {
			log.Fatal(err.Error())
		}
	}(mongoClient, ctx)

	basePath := server.Group("/api").Group("/v1")
	userController.RegisterUserRoutes(basePath)
	authController.RegisterAuthRoutes(basePath)

	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowCredentials: false,
		MaxAge:           int(12 * time.Hour),
	})
	handler := corsConfig.Handler(server)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
