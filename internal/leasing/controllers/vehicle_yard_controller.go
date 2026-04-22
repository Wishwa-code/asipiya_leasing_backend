package controllers

import (
	"garment-management-backend/internal/leasing/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VehicleYardController struct {
	DB *gorm.DB
}

func (ctrl *VehicleYardController) Index(c *gin.Context) {
	var records []models.VehicleYard
	if err := ctrl.DB.Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vehicle yards"})
		return
	}
	c.JSON(http.StatusOK, records)
}

func (ctrl *VehicleYardController) Store(c *gin.Context) {
	var req models.VehicleYard

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vehicle yard"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func (ctrl *VehicleYardController) Update(c *gin.Context) {
	id := c.Param("id")
	var record models.VehicleYard
	if err := ctrl.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle yard not found"})
		return
	}

	var req models.VehicleYard
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = record.ID
	req.CreatedAt = record.CreatedAt

	if err := ctrl.DB.Save(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle yard"})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (ctrl *VehicleYardController) Destroy(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&models.VehicleYard{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vehicle yard"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vehicle yard deleted successfully"})
}
