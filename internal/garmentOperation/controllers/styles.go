package controllers

import (
	// "fmt"
	"net/http"
	// "time"
	"garment-management-backend/internal/garmentOperation/operationModels"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StyleController struct {
	DB *gorm.DB
}

// Store: Create a new Style ➕
func (ctrl *StyleController) Store(c *gin.Context) {
	var input operationModels.StyleRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	style := operationModels.Style{
		StyleCode:         input.StyleCode,
		SMV:               input.SMV,
		StyleCategoryCode: input.StyleCategoryCode,
		StyleSerialNo:     input.StyleSerialNo,
	}

	if err := ctrl.DB.Create(&style).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Style code already exists or database error"})
		return
	}

	c.JSON(http.StatusCreated, style)
}

// Index: List all Styles 📋
func (ctrl *StyleController) Index(c *gin.Context) {
	var styles []operationModels.Style
	ctrl.DB.Find(&styles)
	c.JSON(http.StatusOK, styles)
}

// Show: Get single Style 🔍
func (ctrl *StyleController) Show(c *gin.Context) {
	var style operationModels.Style
	if err := ctrl.DB.First(&style, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Style not found"})
		return
	}
	c.JSON(http.StatusOK, style)
}

// Update: Update Style 📝
func (ctrl *StyleController) Update(c *gin.Context) {
	var style operationModels.Style
	if err := ctrl.DB.First(&style, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Style not found"})
		return
	}

	var input operationModels.StyleRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Model(&style).Updates(operationModels.Style{
		StyleCode:         input.StyleCode,
		SMV:               input.SMV,
		StyleCategoryCode: input.StyleCategoryCode,
		StyleSerialNo:     input.StyleSerialNo,
	})

	c.JSON(http.StatusOK, style)
}

// Destroy: Delete Style 🗑️
func (ctrl *StyleController) Destroy(c *gin.Context) {
	if err := ctrl.DB.Delete(&operationModels.Style{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Deletion failed"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
