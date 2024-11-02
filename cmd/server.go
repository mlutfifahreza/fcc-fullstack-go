package cmd

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"github.com/mlutfifahreza/fcc-fullstack-go/api"
	"github.com/mlutfifahreza/fcc-fullstack-go/internal/product_db"
)

type AppDependencies struct {
	productDB *product_db.Database
}

func RunServer() {
	// Initialize Fiber app
	app := fiber.New()

	// Initialize Dependencies
	dsn := "postgres://username:password@localhost:5432/product_db"
	productDB, err := product_db.NewDatabase(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer productDB.Close()

	dep := &AppDependencies{productDB: productDB}

	// Setup routes
	setupRoute(app, dep)

	// Run
	log.Fatal(app.Listen(fmt.Sprintf(":%d", 8000)))
}

func setupRoute(app *fiber.App, dep *AppDependencies) {
	app.Get("/", api.HandleHome())
	app.Get("/ping", api.HandlePing())

	app.Post("/products", api.HandleCreateProduct(dep.productDB))
	app.Patch("/products", api.HandleUpdateProduct(dep.productDB))
	app.Get("/products/:id", api.HandleGetProduct(dep.productDB))
	app.Delete("/products/:id", api.HandleDeleteProduct(dep.productDB))

	log.Info("Routes are setup")
}
