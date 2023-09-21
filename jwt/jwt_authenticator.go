package jwt

import (
	"greens-basket/utils"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gofiber/fiber/v2"
)

func extractToken(c *fiber.Ctx) (string, error) {

	// infoLog := c.Locals(utils.InfoLogger).(*log.Logger)

	// infoLog.Println("[JWT_Authenticator] extractToken:")

	// Get the "Authorization" header value from the request
	token := c.Get("Authorization")

	// Check if the header value is empty or doesn't start with "Bearer "
	if token == "" || len(token) < 7 || token[:7] != "Bearer " {
		return "", c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid or missing token",
		})
	}

	// Extract the token without "Bearer "
	return token[7:], nil
}

func (jf *JWTFactory) Authenticate(scp string, r utils.Role) fiber.Handler {

	return func(c *fiber.Ctx) error {

		tokenString, err := extractToken(c)

		if err != nil {
			//errLog.Println("[JWT_Authenticator] Protected: error = " + err.Error())
			return err
		}

		//	infoLog.Println("[JWT_Authenticator] Protected token = " + tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jf.Config.TokenSecretKey), nil
		})

		if err != nil {
			//	errLog.Println("[JWT_Authenticator] Protected: error = " + err.Error())
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"status": "error", "message": "Malformed JWT or " + err.Error(), "data": nil})

		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid claims in token",
			})
		}

		// verify scope
		scope := claims[utils.Scope]

		if scope != scp {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid claims in token",
			})
		}

		// verify role
		role := claims[utils.Roles]

		if role != string(r) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid claims in token",
			})
		}

		c.Locals(utils.JWTToken, token)

		return c.Next()
	}

}
