package handler

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"internal-backend/utils"
	"time"
)

func Login(ctx *fiber.Ctx) error {
	var (
		err    error
		secret string //jwt secret from env
	)
	// Represent username and password
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	var input LoginInput

	// Try to parse user/pass from body, return unauthorized otherwise
	if err = ctx.BodyParser(&input); err != nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	// Check user/pass validity vs DB
	// TODO: implement DB interaction
	// FIXME: remove dummy auth
	if input.Identity != "user" || input.Password != "pass" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	// Create new token with claims
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["identity"] = input.Identity
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()

	// Read secret from env
	secret, err = utils.ReadEnv("SECRET")
	if err != nil {
		// TODO: implement error logging
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	// Sign token
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		// TODO: implement error logging
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	// Since all is ok, return json with signed token
	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Success login",
		"data":    t,
	})
}
