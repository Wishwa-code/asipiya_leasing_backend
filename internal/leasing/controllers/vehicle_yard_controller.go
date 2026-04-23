package controllers

import (
	"fmt"
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
	var payload struct {
		YardName      string      `json:"yard_name"`
		Address       string      `json:"address"`
		Province      string      `json:"province"`
		District      string      `json:"district"`
		DSD           string      `json:"dsd"`
		YardContactNo string      `json:"yard_contact_no"`
		ContactPerson string      `json:"contact_person"`
		MobileNo      string      `json:"mobile_no"`
		Status        interface{} `json:"status"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statusStr := "Active"
	if payload.Status != nil {
		statusStr = fmt.Sprintf("%v", payload.Status)
		if statusStr == "1" || statusStr == "true" {
			statusStr = "Active"
		} else if statusStr == "0" || statusStr == "false" {
			statusStr = "Inactive"
		}
	}

	req := models.VehicleYard{
		YardName:      payload.YardName,
		Address:       payload.Address,
		Province:      payload.Province,
		District:      payload.District,
		DSD:           payload.DSD,
		YardContactNo: payload.YardContactNo,
		ContactPerson: payload.ContactPerson,
		MobileNo:      payload.MobileNo,
		Status:        statusStr,
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

	var payload struct {
		YardName      string      `json:"yard_name"`
		Address       string      `json:"address"`
		Province      string      `json:"province"`
		District      string      `json:"district"`
		DSD           string      `json:"dsd"`
		YardContactNo string      `json:"yard_contact_no"`
		ContactPerson string      `json:"contact_person"`
		MobileNo      string      `json:"mobile_no"`
		Status        interface{} `json:"status"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statusStr := record.Status
	if payload.Status != nil {
		statusStr = fmt.Sprintf("%v", payload.Status)
		if statusStr == "1" || statusStr == "true" {
			statusStr = "Active"
		} else if statusStr == "0" || statusStr == "false" {
			statusStr = "Inactive"
		}
	}

	record.YardName = payload.YardName
	record.Address = payload.Address
	record.Province = payload.Province
	record.District = payload.District
	record.DSD = payload.DSD
	record.YardContactNo = payload.YardContactNo
	record.ContactPerson = payload.ContactPerson
	record.MobileNo = payload.MobileNo
	record.Status = statusStr

	if err := ctrl.DB.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle yard"})
		return
	}

	c.JSON(http.StatusOK, record)
}

func (ctrl *VehicleYardController) Destroy(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&models.VehicleYard{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vehicle yard"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vehicle yard deleted successfully"})
}
