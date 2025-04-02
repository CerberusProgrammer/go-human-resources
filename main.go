package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"

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

	// After creating the template engine...
	engine := html.New(templatesPath, ".html")
	engine.Reload(true) // Enable reloading for development

	// Add helper functions to handle string operations safely
	engine.AddFunc("sub", func(a, b int) int {
		return a - b
	})

	engine.AddFunc("ge", func(a, b int) bool {
		return a >= b
	})

	engine.AddFunc("safe", func(s string) template.HTML {
		return template.HTML(s)
	})

	// Force load templates in the correct order
	templatePaths := []string{
		"layout/main.layout.html",
		"index.html",
		"pages/error.html",
		"pages/employee/employee_view.html",
		"pages/employee/employee_edit_view.html",
		"pages/employee/employee_create_view.html",
		"pages/employee/employees_view.html",
		"pages/employee/partials/employee_partial_data.html",
		"pages/employee/partials/employees_partial_list.html",
	}

	for _, path := range templatePaths {
		fullPath := filepath.Join(templatesPath, path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			log.Printf("Warning: Template file not found: %s", fullPath)
		} else {
			log.Printf("Loading template: %s", path)
			engine.Load()
		}
	}
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
