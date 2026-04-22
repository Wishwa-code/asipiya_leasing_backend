package controllers

import (
	"garment-management-backend/internal/leasing/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SeizerController struct {
	DB *gorm.DB
}

// Index handles GET /api/v1/seizers
func (ctrl *SeizerController) Index(c *gin.Context) {
	var seizers []models.Seizer
	if err := ctrl.DB.Find(&seizers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seizers"})
		return
	}
	c.JSON(http.StatusOK, seizers)
}

// Store handles POST /api/v1/seizers
func (ctrl *SeizerController) Store(c *gin.Context) {
	var req models.Seizer

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create seizer"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// Update handles PUT /api/v1/seizers/:id
func (ctrl *SeizerController) Update(c *gin.Context) {
	id := c.Param("id")
	var seizer models.Seizer
	if err := ctrl.DB.First(&seizer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seizer not found"})
		return
	}

	var req models.Seizer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = seizer.ID
	req.CreatedAt = seizer.CreatedAt

	if err := ctrl.DB.Save(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seizer"})
		return
	}

	c.JSON(http.StatusOK, req)
}

// Destroy handles DELETE /api/v1/seizers/:id
func (ctrl *SeizerController) Destroy(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&models.Seizer{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete seizer"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Seizer deleted successfully"})
}
