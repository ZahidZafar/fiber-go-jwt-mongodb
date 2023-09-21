package controller

import (
	"context"
	"greens-basket/data"
	"greens-basket/repositories"
	"greens-basket/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type UserController struct {
	UserRepo *repositories.UserRepository
}

func (uc *UserController) Profile(c *fiber.Ctx) error {
	token := c.Locals(utils.JWTToken).(*jwt.Token)

	claims, _ := token.Claims.(jwt.MapClaims)

	subject := claims[utils.Subject].(string)

	// Find one document that matches the filter
	result, err := uc.UserRepo.GetByID(context.Background(), subject)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found!",
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// Update Profile
func (uc *UserController) UpdateProfile(c *fiber.Ctx) error {
	var user data.AppUser
	if err := c.BodyParser(&user); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	token := c.Locals(utils.JWTToken).(*jwt.Token)

	claims, _ := token.Claims.(jwt.MapClaims)

	subject := claims[utils.Subject].(string)

	var result = &data.AppUser{}

	err := uc.UserRepo.UpdateByID(context.Background(), subject, bson.M{"n": user.Name, "b": user.Balance}, result)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})

	}

	return c.Status(fiber.StatusOK).JSON(result)
}
