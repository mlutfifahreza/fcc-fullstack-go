package util

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ParseAndValidateRequest(ctx *fiber.Ctx, req any) error {
	err := ctx.BodyParser(req)
	if err != nil {
		return err
	}

	err = validate.Struct(req)
	if err != nil {
		return err
	}

	return nil
}
