package operationModels

import "gorm.io/gorm"

type Employee struct {
    gorm.Model
    Name         string `gorm:"not null" json:"name"`
    Address      string `json:"address"`
    EmployeeCode string `gorm:"uniqueIndex;not null" json:"employee_code"`
    MobileNumber string `json:"mobile_number"`
    EfpNumber    string `json:"efp_number"`
}

// EmployeeRequest handles validation for incoming JSON
type EmployeeRequest struct {
    Name         string `json:"name" binding:"required"`
    Address      string `json:"address" binding:"required"`
    EmployeeCode string `json:"employee_code" binding:"required"`
    MobileNumber string `json:"mobile_number"`
    EfpNumber    string `json:"efp_number"`
}