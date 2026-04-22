package controllers

import (
	"fmt"
	"garment-management-backend/internal/leasing/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SupplierController struct {
	DB *gorm.DB
}

// Index handles GET /api/suppliers
func (ctrl *SupplierController) Index(c *gin.Context) {
	var suppliers []models.Supplier
	if err := ctrl.DB.Find(&suppliers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch suppliers"})
		return
	}
	c.JSON(http.StatusOK, suppliers)
}

// Store handles POST /api/suppliers
func (ctrl *SupplierController) Store(c *gin.Context) {
	var req struct {
		Name         string      `json:"name"`
		NIC          string      `json:"nic"`
		ContactNo    string      `json:"contact_no"`
		Occupation   string      `json:"occupation"`
		Income       interface{} `json:"income"` // Handle both string and number from frontend
		NameInCheque string      `json:"name_in_cheque"`
		Latitude     string      `json:"latitude"`
		Longitude    string      `json:"longitude"`
		Address      string      `json:"address"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse numeric values
	income, _ := strconv.ParseFloat(fmt.Sprintf("%v", req.Income), 64)
	lat, _ := strconv.ParseFloat(req.Latitude, 64)
	lng, _ := strconv.ParseFloat(req.Longitude, 64)

	supplier := models.Supplier{
		Name:         req.Name,
		NIC:          req.NIC,
		ContactNo:    req.ContactNo,
		Occupation:   req.Occupation,
		Income:       income,
		NameInCheque: req.NameInCheque,
		Latitude:     lat,
		Longitude:    lng,
		Address:      req.Address,
	}

	if err := ctrl.DB.Create(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create supplier"})
		return
	}

	c.JSON(http.StatusCreated, supplier)
}

// Update handles PUT /api/suppliers/:id
func (ctrl *SupplierController) Update(c *gin.Context) {
	id := c.Param("id")
	var supplier models.Supplier
	if err := ctrl.DB.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		return
	}

	var req struct {
		Name         string      `json:"name"`
		NIC          string      `json:"nic"`
		ContactNo    string      `json:"contact_no"`
		Occupation   string      `json:"occupation"`
		Income       interface{} `json:"income"`
		NameInCheque string      `json:"name_in_cheque"`
		Latitude     string      `json:"latitude"`
		Longitude    string      `json:"longitude"`
		Address      string      `json:"address"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	income, _ := strconv.ParseFloat(fmt.Sprintf("%v", req.Income), 64)
	lat, _ := strconv.ParseFloat(req.Latitude, 64)
	lng, _ := strconv.ParseFloat(req.Longitude, 64)

	supplier.Name = req.Name
	supplier.NIC = req.NIC
	supplier.ContactNo = req.ContactNo
	supplier.Occupation = req.Occupation
	supplier.Income = income
	supplier.NameInCheque = req.NameInCheque
	supplier.Latitude = lat
	supplier.Longitude = lng
	supplier.Address = req.Address

	if err := ctrl.DB.Save(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update supplier"})
		return
	}

	c.JSON(http.StatusOK, supplier)
}

// Destroy handles DELETE /api/suppliers/:id
func (ctrl *SupplierController) Destroy(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&models.Supplier{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete supplier"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Supplier deleted successfully"})
}
