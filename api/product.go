package api

import (
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/mlutfifahreza/fcc-fullstack-go/internal/product_db"
	"github.com/mlutfifahreza/fcc-fullstack-go/pkg/util"
)

type NewProduct struct {
	Name string `json:"name" validate:"required"`
}

func HandleCreateProduct(productDb *product_db.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var newProduct NewProduct

		if err := util.ParseAndValidateRequest(ctx, &newProduct); err != nil {
			return BadRequestResponse(ctx, err)
		}

		product := product_db.Product{
			ID:   uuid.New().String(),
			Name: newProduct.Name,
		}

		if err := productDb.CreateProduct(ctx.UserContext(), &product); err != nil {
			return InternalErrorResponse(ctx, err)
		}

		return SuccessResponse(ctx, product)
	}
}

func HandleGetProduct(productDb *product_db.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")

		product, err := productDb.GetProduct(ctx.UserContext(), id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return NotFoundResponse(ctx)
			}
			return InternalErrorResponse(ctx, err)
		}

		return SuccessResponse(ctx, product)
	}
}

func HandleGetProductList(productDb *product_db.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		page, _ := strconv.Atoi(ctx.Query("page", "0"))
		size, _ := strconv.Atoi(ctx.Query("size", "0"))

		if size == 0 {
			size = 10
		}

		if page == 0 {
			page = 1
		}

		filter := product_db.GetProductListFilter{
			Limit:  size,
			Offset: (page - 1) * size,
		}

		products, err := productDb.GetProductList(ctx.UserContext(), filter)
		if err != nil {
			return InternalErrorResponse(ctx, err)
		}

		return SuccessResponse(ctx, products)
	}
}

func HandleUpdateProduct(productDb *product_db.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var updatedProduct product_db.Product

		if err := util.ParseAndValidateRequest(ctx, &updatedProduct); err != nil {
			return BadRequestResponse(ctx, err)
		}

		updatedProduct.UpdatedAt = time.Now()

		if err := productDb.UpdateProduct(ctx.UserContext(), &updatedProduct); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return NotFoundResponse(ctx)
			}
			return InternalErrorResponse(ctx, err)
		}

		return SuccessResponse(ctx, nil)
	}
}

func HandleDeleteProduct(productDb *product_db.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")

		if err := productDb.DeleteProduct(ctx.UserContext(), id); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return NotFoundResponse(ctx)
			}
			return InternalErrorResponse(ctx, err)
		}

		return SuccessResponse(ctx, fiber.Map{"id": id})
	}
}
