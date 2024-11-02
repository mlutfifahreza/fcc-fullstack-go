package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Code    string `json:"code,omitempty"`
}

func SuccessResponse(ctx *fiber.Ctx, data any) error {
	return ctx.Status(200).JSON(data)
}

func BadRequestResponse(ctx *fiber.Ctx, err error) error {
	return ctx.Status(400).JSON(ErrorResponse{
		Message: "Invalid Request",
		Error:   err.Error(),
	})
}

func NotFoundResponse(ctx *fiber.Ctx) error {
	return ctx.Status(404).JSON(fiber.Map{"message": "Not found"})
}

func InternalErrorResponse(ctx *fiber.Ctx, err error) error {
	log.Errorw("InternalErrorResponse", fiber.Map{
		"URL":       ctx.OriginalURL(),
		"Method":    ctx.Method(),
		"Queries":   ctx.Queries(),
		"AllParams": ctx.AllParams(),
		"Body":      ctx.Body(),
		"Error":     err.Error(),
	})
	return ctx.Status(500).JSON(fiber.Map{"message": "Internal Server Error"})
}
