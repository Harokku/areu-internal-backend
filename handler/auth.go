package handler

import (
	"github.com/gofiber/fiber/v2"
	"internal-backend/auth"
	"internal-backend/database"
)

func Login(ctx *fiber.Ctx) error {
	var (
		err    error
		u      database.User          //User retrieved from db
		claims map[string]interface{} //Claims to be added to jwt
		t      string                 //Signed token
	)
	// Represent username and password
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	var input LoginInput

	// Try to parse user/pass from body, return unauthorized otherwise
	if err = ctx.BodyParser(&input); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	// -------------------------
	// Check user/pass vs db
	// -------------------------

	// Get user's data from db
	err = u.Get(input.Identity)
	if err != nil {
		// TODO: implement error logging
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	// Check if password match to db
	if !auth.ComparePassword(u.Password, input.Password) {
		// TODO: implement error logging
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	// Claims to be added to JWT
	claims = make(map[string]interface{})
	claims["identity"] = input.Identity
	claims["name"] = u.Name
	claims["surname"] = u.Surname

	// Sign token
	t, err = auth.CreateJWT(claims)
	if err != nil {
		//TODO: implement error logging
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	// All ok, return signed token
	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Success login",
		"data":    t,
	})
}
