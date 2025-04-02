package employee

import (
	"errors"
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

	// Basic Information
	employee.FirstName = c.FormValue("first_name")
	employee.LastName = c.FormValue("last_name")
	employee.Email = c.FormValue("email")
	employee.Phone = c.FormValue("phone")
	employee.Position = c.FormValue("position")
	employee.Department = c.FormValue("department")
	employee.EmploymentType = c.FormValue("employment_type")
	employee.ProfilePicture = c.FormValue("profile_picture")

	// Validate required fields
	if employee.FirstName == "" || employee.LastName == "" || employee.Email == "" || c.FormValue("ssn") == "" {
		return c.Status(fiber.StatusBadRequest).Render("pages/error", fiber.Map{
			"Message":         "First name, last name, email, and SSN are required",
			"CurrentTemplate": "error",
		})
	}

	// Check if email already exists
	var existingEmployee Employee
	if err := h.DB.Where("email = ?", employee.Email).First(&existingEmployee).Error; err == nil {
		// Email exists
		return c.Status(fiber.StatusBadRequest).Render("pages/error", fiber.Map{
			"Message":         "An employee with this email address already exists",
			"CurrentTemplate": "error",
		})
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Database error
		log.Printf("Error checking email uniqueness: %v", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/error", fiber.Map{
			"Message":         "Failed to validate email uniqueness",
			"CurrentTemplate": "error",
		})
	}

	// Personal Information
	employee.SSN = c.FormValue("ssn")
	employee.Address = c.FormValue("address")
	employee.City = c.FormValue("city")
	employee.State = c.FormValue("state")
	employee.PostalCode = c.FormValue("postal_code")
	employee.Country = c.FormValue("country")

	// Financial Information
	employee.BankName = c.FormValue("bank_name")
	employee.BankAccount = c.FormValue("bank_account")
	employee.BankRouting = c.FormValue("bank_routing")
	employee.SalaryCurrency = c.FormValue("salary_currency")
	employee.SalaryFrequency = c.FormValue("salary_frequency")

	// Emergency Contact
	employee.EmergencyContactName = c.FormValue("emergency_contact_name")
	employee.EmergencyContactPhone = c.FormValue("emergency_contact_phone")
	employee.EmergencyContactEmail = c.FormValue("emergency_contact_email")

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

	// Parse date of birth from form input
	dobStr := c.FormValue("date_of_birth")
	if dobStr != "" {
		dob, err := time.Parse("2006-01-02", dobStr)
		if err == nil {
			employee.DateOfBirth = dob
		} else {
			log.Printf("Error parsing date of birth: %v", err)
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

	// Update basic information
	employee.FirstName = c.FormValue("first_name")
	employee.LastName = c.FormValue("last_name")
	employee.Email = c.FormValue("email")
	employee.Phone = c.FormValue("phone")
	employee.Position = c.FormValue("position")
	employee.Department = c.FormValue("department")
	employee.EmploymentType = c.FormValue("employment_type")
	employee.ProfilePicture = c.FormValue("profile_picture")

	// Validate required fields
	if employee.FirstName == "" || employee.LastName == "" || employee.Email == "" || c.FormValue("ssn") == "" {
		return c.Status(fiber.StatusBadRequest).Render("pages/error", fiber.Map{
			"Message":         "First name, last name, email, and SSN are required",
			"CurrentTemplate": "error",
		})
	}

	// Personal Information
	employee.SSN = c.FormValue("ssn")
	employee.Address = c.FormValue("address")
	employee.City = c.FormValue("city")
	employee.State = c.FormValue("state")
	employee.PostalCode = c.FormValue("postal_code")
	employee.Country = c.FormValue("country")

	// Financial Information
	employee.BankName = c.FormValue("bank_name")
	employee.BankAccount = c.FormValue("bank_account")
	employee.BankRouting = c.FormValue("bank_routing")
	employee.SalaryCurrency = c.FormValue("salary_currency")
	employee.SalaryFrequency = c.FormValue("salary_frequency")

	// Emergency Contact
	employee.EmergencyContactName = c.FormValue("emergency_contact_name")
	employee.EmergencyContactPhone = c.FormValue("emergency_contact_phone")
	employee.EmergencyContactEmail = c.FormValue("emergency_contact_email")

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

	// Parse date of birth from form input
	dobStr := c.FormValue("date_of_birth")
	if dobStr != "" {
		dob, err := time.Parse("2006-01-02", dobStr)
		if err == nil {
			employee.DateOfBirth = dob
		} else {
			log.Printf("Error parsing date of birth: %v", err)
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

	// Try to find the employee first to make sure it exists
	result := h.DB.First(&employee, id)
	if result.Error != nil {
		log.Printf("Error finding employee for deletion: %v", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).Render("pages/error", fiber.Map{
				"Message":         "Employee not found",
				"CurrentTemplate": "error",
			})
		}
		return c.Status(fiber.StatusInternalServerError).Render("pages/error", fiber.Map{
			"Message":         "Failed to delete employee",
			"CurrentTemplate": "error",
		})
	}

	// Actually delete the employee
	result = h.DB.Delete(&employee)
	if result.Error != nil {
		log.Printf("Error deleting employee: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).Render("pages/error", fiber.Map{
			"Message":         "Failed to delete employee",
			"CurrentTemplate": "error",
		})
	}

	// If HTMX request, return an empty response with 200 status
	if c.Get("HX-Request") == "true" {
		return c.SendString("")
	}

	// Regular form submission redirects to employee list
	return c.Redirect("/employees")
}
