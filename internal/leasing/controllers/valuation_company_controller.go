package controllers

import (
	"garment-management-backend/internal/leasing/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ValuationCompanyController struct {
	DB *gorm.DB
}

func (ctrl *ValuationCompanyController) Index(c *gin.Context) {
	var records []models.ValuationCompany
	if err := ctrl.DB.Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch valuation companies"})
		return
	}
	c.JSON(http.StatusOK, records)
}

func (ctrl *ValuationCompanyController) Store(c *gin.Context) {
	var req models.ValuationCompany

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create valuation company"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func (ctrl *ValuationCompanyController) Update(c *gin.Context) {
	id := c.Param("id")
	var record models.ValuationCompany
	if err := ctrl.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Valuation company not found"})
		return
	}

	var req models.ValuationCompany
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = record.ID
	req.CreatedAt = record.CreatedAt

	if err := ctrl.DB.Save(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update valuation company"})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (ctrl *ValuationCompanyController) Destroy(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&models.ValuationCompany{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete valuation company"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Valuation company deleted successfully"})
}
