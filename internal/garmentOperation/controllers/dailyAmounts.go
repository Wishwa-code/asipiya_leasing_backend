package controllers

import (
	"garment-management-backend/internal/garmentOperation/operationModels"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OperationDataController struct {
	DB *gorm.DB
}

func (ctrl *OperationDataController) GetDailyAmounts(c *gin.Context) {
	var amounts []operationModels.DailyAmount

	// Highlight: Using Preload to include related Style and Employee data
	result := ctrl.DB.
		Preload("Style").
		Preload("Employee").
		Order("created_at desc").
		Find(&amounts)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch operation data"})
		return
	}

	c.JSON(http.StatusOK, amounts)
}
