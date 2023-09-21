package main

import (
	"context"
	"greens-basket/utils"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	controller "greens-basket/controllers"
	"greens-basket/jwt"
	"greens-basket/repositories"
)

type Application struct {
	Fiber  *fiber.App
	Config *utils.Config
	Client *mongo.Client
}

func (app *Application) Init() {
	// // Create logger for writing information and error messages.

	// app.Fiber.Use(func(c *fiber.Ctx) error {
	// 	c.Locals(utils.MongoClient, app.Client) // Store the client in the context
	// 	return c.Next()
	// })
	app.registerRoutes()
}

func (app *Application) registerRoutes() {
	// Register handler functions.

	jFactory := &jwt.JWTFactory{
		Config: app.Config,
	}

	db := app.Client.Database(utils.Database)
	userRepo := repositories.NewUserRepository(db)

	authController := &controller.AuthController{
		JFactor:  jFactory,
		UserRepo: userRepo,
	}

	userController := &controller.UserController{
		UserRepo: userRepo,
	}

	authGroup := app.Fiber.Group("/auth")
	authGroup.Post("/verify-pincode", jFactory.Authenticate(utils.TempToken, ""), authController.VerifyPinCode)
	authGroup.Post("/authenticate", authController.Authenticate)

	userGroup := app.Fiber.Group("/user")
	userGroup.Get("/", jFactory.Authenticate(utils.AccessToken, utils.User), userController.Profile)
	userGroup.Post("/", jFactory.Authenticate(utils.AccessToken, utils.User), userController.UpdateProfile)
	//userGroup.Use(jwt.Protected(utils.AccessToken, utils.User)) // protecting group route with token and role
}

func main() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	app := fiber.New()
	//fiber.Use(SetCommonHeaders())

	//serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.DBSource) //.SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	ap := Application{
		Fiber:  app,
		Config: &config,
		Client: client,
	}

	ap.Init()

	if err != nil {
		panic(err)
	}

	e := app.Listen(":3000")
	if e != nil {
		panic(e)
	}
}
