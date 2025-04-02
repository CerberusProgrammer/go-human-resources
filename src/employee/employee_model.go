package employee

import (
	"time"

	"gorm.io/gorm"
)

// Employee represents an employee in the system
type Employee struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	FirstName  string         `json:"first_name" form:"first_name" gorm:"not null"`
	LastName   string         `json:"last_name" form:"last_name" gorm:"not null"`
	Email      string         `json:"email" form:"email" gorm:"unique;not null"`
	Phone      string         `json:"phone" form:"phone"`
	Position   string         `json:"position" form:"position"`
	Department string         `json:"department" form:"department"`
	HireDate   time.Time      `json:"hire_date" form:"hire_date"`
	Salary     float64        `json:"salary" form:"salary"`
	IsActive   bool           `json:"is_active" form:"is_active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName sets the table name for the Employee model
func (Employee) TableName() string {
	return "employees"
}
