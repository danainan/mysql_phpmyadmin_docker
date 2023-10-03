package handler

import (
	"auth/config"
	"auth/connection"
	"auth/models"
	"fmt"
	"log"
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
	return ctx.JSON(fiber.Map{"status": "success", "message": "Api Connected", "data": nil})
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

	secret := []byte(config.Config("SECRET"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["role"] = "admin"
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	t, err := token.SignedString(secret)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.Status(200).JSON(fiber.Map{"status": "ok", "message": "Login success", "token": t})

}

func GetAllUser(ctx *fiber.Ctx) error {
	db := connection.DBConn
	var users []models.User

	db.Find(&users)

	return ctx.Status(200).JSON(fiber.Map{"status": "ok", "message": "User Read Success", "users": users})
}

func GetById(ctx *fiber.Ctx) error {
	id := ctx.Locals("id").(int)

	fmt.Println("id =>>>>>>>> ", id)
	var user models.User

	db := connection.DBConn

	db.First(&user, id)
	return ctx.Status(200).JSON(fiber.Map{"status": "ok", "message": "User Read Success", "users": user})
}
