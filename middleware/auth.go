package middleware

import (
	"auth/config"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func jwtError(ctx *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return ctx.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

func JWTAuthen() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		secret := []byte(config.Config("SECRET"))

		authHeader := ctx.Get(fiber.HeaderAuthorization)
		authToken := ""
		if authHeader != "" {
			authToken = strings.Split(authHeader, " ")[1]
		}

		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Invalid signing method %v", token.Header["alg"])
			}
			fmt.Println(authToken)

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return nil, fmt.Errorf("Failed to parse claims")
			}
			idClaim, ok := claims["id"].(float64)
			id := int(idClaim)
			ctx.Locals("id", id)
			fmt.Println("id middleware =>>>>> ", id)
			return secret, nil
		})
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		if !token.Valid {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token is not valid",
			})
		}

		return ctx.Next()
	}

}
