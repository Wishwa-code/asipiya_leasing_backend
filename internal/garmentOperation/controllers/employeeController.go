package controllers

import (
	"garment-management-backend/internal/garmentOperation/operationModels"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type EmployeeController struct {
	DB *gorm.DB
}

// Store: Create a new Employee ➕
func (ctrl *EmployeeController) Store(c *gin.Context) {
	var input operationModels.EmployeeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee := operationModels.Employee{
		Name:      input.Name,
		EfpNumber: input.EfpNumber,
	}

	if err := ctrl.DB.Create(&employee).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Employee code already exists or database error"})
		return
	}

	c.JSON(http.StatusCreated, employee)
}

// Index: List all Employees 📋
func (ctrl *EmployeeController) Index(c *gin.Context) {
	var employees []operationModels.Employee
	ctrl.DB.Find(&employees)
	c.JSON(http.StatusOK, employees)
}

// Show: Get single Employee 🔍
func (ctrl *EmployeeController) Show(c *gin.Context) {
	var employee operationModels.Employee
	if err := ctrl.DB.First(&employee, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	c.JSON(http.StatusOK, employee)
}

// Update: Update Employee 📝
func (ctrl *EmployeeController) Update(c *gin.Context) {
	var employee operationModels.Employee
	if err := ctrl.DB.First(&employee, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	var input operationModels.EmployeeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Model(&employee).Updates(operationModels.Employee{
		Name:      input.Name,
		EfpNumber: input.EfpNumber,
	})

	c.JSON(http.StatusOK, employee)
}

// Destroy: Delete Employee 🗑️
func (ctrl *EmployeeController) Destroy(c *gin.Context) {
	if err := ctrl.DB.Delete(&operationModels.Employee{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Deletion failed"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
