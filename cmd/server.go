package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/joho/godotenv"

	"github.com/mlutfifahreza/fcc-fullstack-go/api"
	"github.com/mlutfifahreza/fcc-fullstack-go/internal/product_db"
)

type AppDependencies struct {
	productDB *product_db.Database
}

type EnvConfig struct {
	// App
	ENV  string
	PORT int

	// Database
	PRODUCT_DB_HOST     string
	PRODUCT_DB_PORT     int
	PRODUCT_DB_USER     string
	PRODUCT_DB_PASSWORD string
	PRODUCT_DB_NAME     string
}

func RunServer() {
	// Initialize Fiber app
	app := fiber.New()

	// Initialize Dependencies
	envConfig := loadEnv()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		envConfig.PRODUCT_DB_USER,
		envConfig.PRODUCT_DB_PASSWORD,
		envConfig.PRODUCT_DB_HOST,
		envConfig.PRODUCT_DB_PORT,
		envConfig.PRODUCT_DB_NAME,
	)
	productDB, err := product_db.NewDatabase(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer productDB.Close()

	dep := &AppDependencies{productDB: productDB}

	// Setup routes
	setupRoute(app, dep)

	// Run
	log.Fatal(app.Listen(fmt.Sprintf(":%d", envConfig.PORT)))
}

func loadEnv() *EnvConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Error loading .env PORT")
	}

	dbPort, err := strconv.Atoi(os.Getenv("PRODUCT_DB_PORT"))
	if err != nil {
		log.Fatalf("Error loading .env PRODUCT_DB_PORT")
	}

	return &EnvConfig{
		ENV:                 os.Getenv("ENV"),
		PORT:                port,
		PRODUCT_DB_HOST:     os.Getenv("PRODUCT_DB_HOST"),
		PRODUCT_DB_PORT:     dbPort,
		PRODUCT_DB_USER:     os.Getenv("PRODUCT_DB_USER"),
		PRODUCT_DB_PASSWORD: os.Getenv("PRODUCT_DB_PASSWORD"),
		PRODUCT_DB_NAME:     os.Getenv("PRODUCT_DB_NAME"),
	}
}

func setupRoute(app *fiber.App, dep *AppDependencies) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin,Content-Type,Accept",
	}))

	app.Get("/", api.HandleHome())
	app.Get("/ping", api.HandlePing())

	app.Post("/products", api.HandleCreateProduct(dep.productDB))
	app.Patch("/products", api.HandleUpdateProduct(dep.productDB))
	app.Get("/products", api.HandleGetProductList(dep.productDB))
	app.Get("/products/:id", api.HandleGetProduct(dep.productDB))
	app.Delete("/products/:id", api.HandleDeleteProduct(dep.productDB))
}
