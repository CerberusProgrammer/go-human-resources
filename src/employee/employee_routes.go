package employee

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes sets up all employee-related routes
func SetupRoutes(app *fiber.App, db *gorm.DB) {
	handler := NewHandler(db)

	// Employee routes
	app.Get("/employees", handler.GetAllEmployees)
	app.Get("/employees/new", handler.ShowCreateForm)
	app.Post("/employees", handler.CreateEmployee)
	app.Get("/employees/:id", handler.GetEmployee)
	app.Get("/employees/:id/edit", handler.ShowEditForm)
	app.Post("/employees/:id", handler.UpdateEmployee)
	app.Delete("/employees/:id", handler.DeleteEmployee)
	app.Post("/employees/:id/delete", handler.DeleteEmployee) // For form submissions
}
