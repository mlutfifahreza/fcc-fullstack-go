package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func HandleHome() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		log.Info(fiber.Map{
			"message": "HandleHome",
			"params":  ctx.AllParams(),
		})

		return ctx.Status(200).JSON(fiber.Map{"data": "hello there"})
	}
}

func HandlePing() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		log.Info(fiber.Map{
			"message": "HandlePing",
			"params":  ctx.AllParams(),
		})

		return ctx.Status(200).SendString("PONG")
	}
}
