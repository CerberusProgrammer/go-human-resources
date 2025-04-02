package employee

import (
	"time"

	"gorm.io/gorm"
)

// Employee represents an employee in the system
type Employee struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	FirstName      string    `json:"first_name" form:"first_name" gorm:"not null"`
	LastName       string    `json:"last_name" form:"last_name" gorm:"not null"`
	Email          string    `json:"email" form:"email" gorm:"unique;not null"`
	Phone          string    `json:"phone" form:"phone"`
	Position       string    `json:"position" form:"position"`
	Department     string    `json:"department" form:"department"`
	HireDate       time.Time `json:"hire_date" form:"hire_date"`
	Salary         float64   `json:"salary" form:"salary"`
	IsActive       bool      `json:"is_active" form:"is_active" gorm:"default:true"`
	ProfilePicture string    `json:"profile_picture" form:"profile_picture"`

	DateOfBirth    time.Time `json:"date_of_birth" form:"date_of_birth"`
	Address        string    `json:"address" form:"address"`
	City           string    `json:"city" form:"city"`
	State          string    `json:"state" form:"state"`
	PostalCode     string    `json:"postal_code" form:"postal_code"`
	Country        string    `json:"country" form:"country"`
	EmploymentType string    `json:"employment_type" form:"employment_type"` // e.g., full-time, part-time, contract

	// Money related fields
	SSN             string `json:"ssn" form:"ssn" gorm:"unique;not null"`    // Social Security Number
	BankAccount     string `json:"bank_account" form:"bank_account"`         // Bank account number
	BankName        string `json:"bank_name" form:"bank_name"`               // Bank name
	BankRouting     string `json:"bank_routing" form:"bank_routing"`         // Bank routing number
	SalaryCurrency  string `json:"salary_currency" form:"salary_currency"`   // Currency of the salary
	SalaryFrequency string `json:"salary_frequency" form:"salary_frequency"` // e.g., monthly, bi-weekly, weekly

	// Emergency contact fields
	EmergencyContactName  string `json:"emergency_contact_name" form:"emergency_contact_name"`
	EmergencyContactPhone string `json:"emergency_contact_phone" form:"emergency_contact_phone"`
	EmergencyContactEmail string `json:"emergency_contact_email" form:"emergency_contact_email"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName sets the table name for the Employee model
func (Employee) TableName() string {
	return "employees"
}
