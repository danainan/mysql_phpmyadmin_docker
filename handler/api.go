package handler

import (
	"auth/connection"
	"auth/models"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func Hello(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{"status": "success", "message": "Hello fiber!!!!", "data": nil})
}

func Register(ctx *fiber.Ctx) error {
	db := connection.DBConn
	user := new(models.User)

	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't register", "data": nil})
	}

	encryptedPassword, err := hashPassword(user.Password)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Password hashing error", "data": nil})
	}
	user.Password = encryptedPassword

	result := db.Create(&user)
	if result.Error != nil {
		log.Println(result.Error)
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Email exist couldn't register", "data": nil})
	}
	return ctx.Status(200).JSON(fiber.Map{"status": "ok", "message": "User Create Success", "data": user.ID})

}

func Login(ctx *fiber.Ctx) error {
	db := connection.DBConn
	var req models.Request

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid request"})
	}

	var user models.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == nil {
			return ctx.Status(401).JSON(fiber.Map{"status": "error", "message": "Email not found"})
		}
	}

	if !checkPasswordHash(req.Password, user.Password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid password", "data": nil})
	}

	secret := []byte("SECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["role"] = "admin"
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	t, err := token.SignedString(secret)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.Status(200).JSON(fiber.Map{"status": "ok", "message": "Login success", "token": t})

}

func GetAllUser(ctx *fiber.Ctx) error {
	db := connection.DBConn
	secret := []byte("SECRET")
	var users []models.User
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
		return secret, nil
	})
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}
	if !token.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token is not valid",
		})
	}

	var user models.User
	if err := db.Where("email = ?", user.Email).First(&user).Error; err != nil {
		if err == nil {
			return ctx.Status(401).JSON(fiber.Map{"status": "error", "message": "Email not found"})
		}
	}

	db.Find(&users)

	return ctx.Status(200).JSON(fiber.Map{"status": "ok", "message": "User Read Success", "users": users})
}
