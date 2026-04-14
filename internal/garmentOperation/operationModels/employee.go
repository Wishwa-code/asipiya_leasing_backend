package operationModels

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	Name      string `gorm:"not null" json:"name"`
	EfpNumber string `json:"efp_number"`
}

// EmployeeRequest handles validation for incoming JSON
type EmployeeRequest struct {
	Name      string `json:"name" binding:"required"`
	EfpNumber string `json:"efp_number" binding:"required"`
}
