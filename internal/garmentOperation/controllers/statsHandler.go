package controllers

import (
	"garment-management-backend/internal/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStatsHandler(c *gin.Context) {
	var count int64
	result := database.DB.Table("products").Count(&count)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_products": count})
}
