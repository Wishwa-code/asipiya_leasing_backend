package controllers

import (
	"garment-management-backend/internal/leasing/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InsuranceCompanyController struct {
	DB *gorm.DB
}

func (ctrl *InsuranceCompanyController) Index(c *gin.Context) {
	var records []models.InsuranceCompany
	if err := ctrl.DB.Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch insurance companies"})
		return
	}
	c.JSON(http.StatusOK, records)
}

func (ctrl *InsuranceCompanyController) Store(c *gin.Context) {
	var req models.InsuranceCompany

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create insurance company"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func (ctrl *InsuranceCompanyController) Update(c *gin.Context) {
	id := c.Param("id")
	var record models.InsuranceCompany
	if err := ctrl.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Insurance company not found"})
		return
	}

	var req models.InsuranceCompany
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = record.ID
	req.CreatedAt = record.CreatedAt

	if err := ctrl.DB.Save(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update insurance company"})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (ctrl *InsuranceCompanyController) Destroy(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&models.InsuranceCompany{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete insurance company"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Insurance company deleted successfully"})
}
