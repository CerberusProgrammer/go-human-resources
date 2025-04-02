package main

import (
	"log"
	"os"

	"go-human-resources/src"
	"go-human-resources/src/employee"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Load configuration
	config := src.LoadConfig()

	// Initialize database
	db, err := src.InitDB(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Debug template directory
	cwd, _ := os.Getwd()
	log.Printf("Current working directory: %s", cwd)

	// Check if templates directory exists
	templatesPath := "./templates"
	if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
		log.Fatalf("Templates directory not found at %s", templatesPath)
	}
	log.Printf("Templates directory exists at %s", templatesPath)

	// Initialize template engine - much simpler setup
	engine := html.New(templatesPath, ".html")
	engine.Reload(true) // Enable reloading for development

	// Create fiber app with the template engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Static files
	app.Static("/static", "./static")

	// Home route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":           "GO Human Resources",
			"CurrentTemplate": "index",
		})
	})

	// Setup employee routes
	employee.SetupRoutes(app, db.DB)

	// Start server
	log.Printf("Server starting on port %s", config.ServerPort)
	log.Fatal(app.Listen(":" + config.ServerPort))
}
