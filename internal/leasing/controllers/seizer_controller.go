package controllers

import (
	"fmt"
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
	var payload struct {
		SeizerType          string      `json:"seizer_type"`
		CompanyName         string      `json:"company_name"`
		CompanyRegistration string      `json:"company_registration"`
		CompanyContactNo    string      `json:"company_contact_no"`
		NIC                 string      `json:"nic"`
		SeizerContactNo     string      `json:"seizer_contact_no"`
		MobileNo            string      `json:"mobile_no"`
		Address             string      `json:"address"`
		Remarks             string      `json:"remarks"`
		Status              interface{} `json:"status"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statusStr := "Active"
	if payload.Status != nil {
		statusStr = fmt.Sprintf("%v", payload.Status)
		if statusStr == "1" {
			statusStr = "Active"
		} else if statusStr == "0" {
			statusStr = "Inactive"
		}
	}

	req := models.Seizer{
		SeizerType:          payload.SeizerType,
		CompanyName:         payload.CompanyName,
		CompanyRegistration: payload.CompanyRegistration,
		CompanyContactNo:    payload.CompanyContactNo,
		NIC:                 payload.NIC,
		SeizerContactNo:     payload.SeizerContactNo,
		MobileNo:            payload.MobileNo,
		Address:             payload.Address,
		Remarks:             payload.Remarks,
		Status:              statusStr,
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

	var payload struct {
		SeizerType          string      `json:"seizer_type"`
		CompanyName         string      `json:"company_name"`
		CompanyRegistration string      `json:"company_registration"`
		CompanyContactNo    string      `json:"company_contact_no"`
		NIC                 string      `json:"nic"`
		SeizerContactNo     string      `json:"seizer_contact_no"`
		MobileNo            string      `json:"mobile_no"`
		Address             string      `json:"address"`
		Remarks             string      `json:"remarks"`
		Status              interface{} `json:"status"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statusStr := seizer.Status
	if payload.Status != nil {
		statusStr = fmt.Sprintf("%v", payload.Status)
		if statusStr == "1" {
			statusStr = "Active"
		} else if statusStr == "0" {
			statusStr = "Inactive"
		}
	}

	seizer.SeizerType = payload.SeizerType
	seizer.CompanyName = payload.CompanyName
	seizer.CompanyRegistration = payload.CompanyRegistration
	seizer.CompanyContactNo = payload.CompanyContactNo
	seizer.NIC = payload.NIC
	seizer.SeizerContactNo = payload.SeizerContactNo
	seizer.MobileNo = payload.MobileNo
	seizer.Address = payload.Address
	seizer.Remarks = payload.Remarks
	seizer.Status = statusStr

	if err := ctrl.DB.Save(&seizer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seizer"})
		return
	}

	c.JSON(http.StatusOK, seizer)
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
