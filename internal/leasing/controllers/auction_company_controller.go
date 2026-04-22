package controllers

import (
	"garment-management-backend/internal/leasing/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuctionCompanyController struct {
	DB *gorm.DB
}

func (ctrl *AuctionCompanyController) Index(c *gin.Context) {
	var records []models.AuctionCompany
	if err := ctrl.DB.Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch auction companies"})
		return
	}
	c.JSON(http.StatusOK, records)
}

func (ctrl *AuctionCompanyController) Store(c *gin.Context) {
	var req models.AuctionCompany

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create auction company"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func (ctrl *AuctionCompanyController) Update(c *gin.Context) {
	id := c.Param("id")
	var record models.AuctionCompany
	if err := ctrl.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Auction company not found"})
		return
	}

	var req models.AuctionCompany
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = record.ID
	req.CreatedAt = record.CreatedAt

	if err := ctrl.DB.Save(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update auction company"})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (ctrl *AuctionCompanyController) Destroy(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&models.AuctionCompany{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete auction company"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Auction company deleted successfully"})
}
