package controllers

import (
	"garment-management-backend/internal/leasing/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IntroducerController struct {
	DB *gorm.DB
}

func (ctrl *IntroducerController) Index(c *gin.Context) {
	var records []models.Introducer
	if err := ctrl.DB.Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch introducers"})
		return
	}
	c.JSON(http.StatusOK, records)
}

func (ctrl *IntroducerController) Store(c *gin.Context) {
	var req models.Introducer

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create introducer"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func (ctrl *IntroducerController) Update(c *gin.Context) {
	id := c.Param("id")
	var record models.Introducer
	if err := ctrl.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Introducer not found"})
		return
	}

	var req models.Introducer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = record.ID
	req.CreatedAt = record.CreatedAt

	if err := ctrl.DB.Save(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update introducer"})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (ctrl *IntroducerController) Destroy(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&models.Introducer{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete introducer"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Introducer deleted successfully"})
}
