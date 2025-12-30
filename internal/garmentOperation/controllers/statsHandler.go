package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "garment-management-backend/internal/database"
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