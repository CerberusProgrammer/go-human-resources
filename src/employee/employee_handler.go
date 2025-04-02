package employee

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Handler contains all handlers for employee routes
type Handler struct {
	DB *gorm.DB
}

// NewHandler creates a new employee handler
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) GetAllEmployees(c *fiber.Ctx) error {
	var employees []Employee

	result := h.DB.Find(&employees)
	if result.Error != nil {
		log.Printf("Error fetching employees: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).Render("pages/error", fiber.Map{
			"Message":         "Failed to retrieve employees",
			"CurrentTemplate": "error",
		})
	}

	return c.Render("pages/employee/employees_view", fiber.Map{
		"Employees":       employees,
		"CurrentTemplate": "employees",
	})
}

func (h *Handler) GetEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	var employee Employee

	result := h.DB.First(&employee, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).Render("pages/error", fiber.Map{
				"Message":         "Employee not found",
				"CurrentTemplate": "error",
			})
		}
		return c.Status(fiber.StatusInternalServerError).Render("pages/error", fiber.Map{
			"Message":         "Failed to retrieve employee",
			"CurrentTemplate": "error",
		})
	}

	return c.Render("pages/employee/employee_view", fiber.Map{
		"Employee":        employee,
		"CurrentTemplate": "employee-view",
	})
}

// ShowEditForm displays the form for editing an employee
func (h *Handler) ShowEditForm(c *fiber.Ctx) error {
	id := c.Params("id")
	var employee Employee

	result := h.DB.First(&employee, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).Render("pages/error", fiber.Map{
				"Message":         "Employee not found",
				"CurrentTemplate": "error",
			})
		}
		return c.Status(fiber.StatusInternalServerError).Render("pages/error", fiber.Map{
			"Message":         "Failed to retrieve employee",
			"CurrentTemplate": "error",
		})
	}

	return c.Render("pages/employee/employee_edit_view", fiber.Map{
		"Employee":        employee,
		"CurrentTemplate": "employee-edit",
	})
}

func (h *Handler) ShowCreateForm(c *fiber.Ctx) error {
	return c.Render("pages/employee/employee_create_view", fiber.Map{
		"CurrentTemplate": "employee-create", // This was missing
	})
}

// CreateEmployee handles the creation of a new employee
func (h *Handler) CreateEmployee(c *fiber.Ctx) error {
	employee := new(Employee)

	// Get form values manually instead of using BodyParser
	employee.FirstName = c.FormValue("first_name")
	employee.LastName = c.FormValue("last_name")
	employee.Email = c.FormValue("email")
	employee.Phone = c.FormValue("phone")
	employee.Position = c.FormValue("position")
	employee.Department = c.FormValue("department")

	// Validate required fields
	if employee.FirstName == "" || employee.LastName == "" || employee.Email == "" {
		return c.Status(fiber.StatusBadRequest).Render("pages/error", fiber.Map{
			"Message":         "First name, last name and email are required",
			"CurrentTemplate": "error",
		})
	}

	// Parse hire date from form input
	hireDateStr := c.FormValue("hire_date")
	if hireDateStr != "" {
		hireDate, err := time.Parse("2006-01-02", hireDateStr)
		if err == nil {
			employee.HireDate = hireDate
		} else {
			log.Printf("Error parsing hire date: %v", err)
		}
	}

	// Parse salary
	salaryStr := c.FormValue("salary")
	if salaryStr != "" {
		salary, err := strconv.ParseFloat(salaryStr, 64)
		if err == nil {
			employee.Salary = salary
		} else {
			log.Printf("Error parsing salary: %v", err)
		}
	}

	// Set isActive based on checkbox
	employee.IsActive = c.FormValue("is_active") == "on"

	// Debug output to see what we're saving
	log.Printf("Creating employee: %+v", employee)

	result := h.DB.Create(employee)
	if result.Error != nil {
		log.Printf("Error creating employee: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).Render("pages/error", fiber.Map{
			"Message":         "Failed to create employee",
			"CurrentTemplate": "error",
		})
	}

	// If HTMX request, return the partial for the new employee
	if c.Get("HX-Request") == "true" {
		return c.Render("pages/employee/partials/employee_partial_data", fiber.Map{
			"Employee": employee,
		})
	}

	// Regular form submission redirects to employee list
	return c.Redirect("/employees")
}

// UpdateEmployee updates an existing employee
func (h *Handler) UpdateEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	var employee Employee

	// First find the employee
	result := h.DB.First(&employee, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).Render("pages/error", fiber.Map{
			"Message":         "Employee not found",
			"CurrentTemplate": "error",
		})
	}

	// Update fields manually from form values
	employee.FirstName = c.FormValue("first_name")
	employee.LastName = c.FormValue("last_name")
	employee.Email = c.FormValue("email")
	employee.Phone = c.FormValue("phone")
	employee.Position = c.FormValue("position")
	employee.Department = c.FormValue("department")

	// Parse hire date from form input
	hireDateStr := c.FormValue("hire_date")
	if hireDateStr != "" {
		hireDate, err := time.Parse("2006-01-02", hireDateStr)
		if err == nil {
			employee.HireDate = hireDate
		}
	}

	// Parse salary
	salaryStr := c.FormValue("salary")
	if salaryStr != "" {
		salary, err := strconv.ParseFloat(salaryStr, 64)
		if err == nil {
			employee.Salary = salary
		}
	}

	// Set isActive based on checkbox
	employee.IsActive = c.FormValue("is_active") == "on"

	// Update the employee
	result = h.DB.Save(&employee)
	if result.Error != nil {
		log.Printf("Error updating employee: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).Render("pages/error", fiber.Map{
			"Message":         "Failed to update employee",
			"CurrentTemplate": "error",
		})
	}

	// If HTMX request, return the partial for the updated employee
	if c.Get("HX-Request") == "true" {
		return c.Render("pages/employee/partials/employee_partial_data", fiber.Map{
			"Employee": employee,
		})
	}

	// Regular form submission redirects to employee list
	return c.Redirect("/employees")
}

// DeleteEmployee deletes an employee
func (h *Handler) DeleteEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	var employee Employee

	result := h.DB.Delete(&employee, id)
	if result.Error != nil {
		log.Printf("Error deleting employee: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).Render("pages/error", fiber.Map{
			"Message": "Failed to delete employee",
		})
	}

	// If HTMX request, return an empty response with 200 status
	if c.Get("HX-Request") == "true" {
		return c.SendString("")
	}

	// Regular form submission redirects to employee list
	return c.Redirect("/employees")
}
