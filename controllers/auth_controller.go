package controller

import (
	"context"
	"greens-basket/data"
	jwt_utils "greens-basket/jwt"
	"greens-basket/repositories"
	"greens-basket/utils"
	"strconv"

	"github.com/golang-jwt/jwt/v4"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type AuthRequest struct {
	Phone string `json:"phone" validate:"required"`
	OTP   string `json:"otp" validate:"required"`
}

type AuthController struct {
	JFactor  *jwt_utils.JWTFactory
	UserRepo *repositories.UserRepository
}

func (ac *AuthController) VerifyPinCode(c *fiber.Ctx) error {
	token := c.Locals(utils.JWTToken).(*jwt.Token)

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid claims in token",
		})
	}

	scope := claims[utils.Scope]

	if scope != utils.TempToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid claims in token",
		})
	}

	var req AuthRequest

	if err := c.BodyParser(&req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user",
		})
	}

	v := utils.EncryptWithSHA512(req.Phone + req.OTP)

	checkSum := claims[utils.Subject]

	if checkSum != v {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid checksum",
		})
	}

	user := data.AppUser{
		Phone: req.Phone,
		Role:  string(utils.User),
	}

	// Define a filter to specify the document you want to retrieve
	filter := bson.M{"ph": user.Phone}

	// Create a variable to store the result
	var up data.UserPrincipal
	var result data.AppUser

	// Find one document that matches the filter
	err := ac.UserRepo.FindOne(context.Background(), filter, &result)

	if err == nil {
		up = data.UserPrincipal{
			Subject: result.ID,
			Roles:   string(utils.User),
		}
	} else {
		e := ac.UserRepo.Insert(context.Background(), &user)
		if e != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": e.Error(),
			})
		}

		up = data.UserPrincipal{
			Subject: user.ID,
			Roles:   string(utils.User),
		}
	}

	err = ac.setTokensInHeader(up, c)
	if err != err {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Account Created!",
	})
}

func (ac *AuthController) setTokensInHeader(up data.UserPrincipal, c *fiber.Ctx) error {
	accessToken, err := ac.JFactor.CreateJWTAccessToken(&up)

	if err != err {
		return err
	}

	refreshToken, err := ac.JFactor.CreateJWTRefreshToken(&up)

	if err != err {
		return err
	}

	// Set access, and refresh token in headers
	c.Set(utils.AccessToken, accessToken)
	c.Set(utils.RefreshToken, refreshToken)

	return nil
}

func (ac *AuthController) Authenticate(c *fiber.Ctx) error {
	auth := new(AuthRequest)
	validate := validator.New()

	err := validate.Struct(auth)
	if err != nil {
		// Validation failed
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			fiber.NewError(fiber.ErrBadRequest.Code, ("Field %s is required.\n"), e.Field())
		}
	}

	if err := c.BodyParser(auth); err != nil {
		return err
	}

	pinCode := strconv.Itoa(utils.RandomNumber())

	sub := utils.EncryptWithSHA512(auth.Phone + pinCode)
	up := data.UserPrincipal{
		Subject: sub,
		Roles:   "",
	}

	token, err := ac.JFactor.CreateJWTTempToken(&up)

	if err != err {
		return fiber.NewError(fiber.ErrBadRequest.Code, err.Error())
	}

	// Set multiple headers
	c.Set(utils.TempToken, token)

	return c.SendString("Pin code sent successfully :" + pinCode)
}
