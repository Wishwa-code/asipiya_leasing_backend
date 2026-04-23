package controllers

import (
	"encoding/json"
	"fmt"
	"garment-management-backend/internal/leasing/models"
	"net/http"
	"strconv"

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
	var payload struct {
		IntroducerType   string      `json:"introducer_type"`
		Name             string      `json:"name"`
		RegistrationNo   string      `json:"registration_no"`
		ContactPerson    string      `json:"contact_person"`
		PrimaryContact   string      `json:"primary_contact"`
		SecondaryContact string      `json:"secondary_contact"`
		Email            string      `json:"email"`
		Address          string      `json:"address"`
		CommissionRate   interface{} `json:"commission_rate"`
		BankDetails      interface{} `json:"bank_details"`
		Remarks          string      `json:"remarks"`
		Status           interface{} `json:"status"`
		CreatedBy        uint        `json:"created_by"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commissionRate, _ := strconv.ParseFloat(fmt.Sprintf("%v", payload.CommissionRate), 64)

	statusStr := "Active"
	if payload.Status != nil {
		statusStr = fmt.Sprintf("%v", payload.Status)
		if statusStr == "1" {
			statusStr = "Active"
		} else if statusStr == "0" {
			statusStr = "Inactive"
		}
	}

	bankDetailsStr := "{}"
	if payload.BankDetails != nil {
		switch v := payload.BankDetails.(type) {
		case string:
			if v != "" {
				bankDetailsStr = v
			}
		default:
			b, _ := json.Marshal(v)
			bankDetailsStr = string(b)
		}
	}

	req := models.Introducer{
		IntroducerType:   payload.IntroducerType,
		Name:             payload.Name,
		RegistrationNo:   payload.RegistrationNo,
		ContactPerson:    payload.ContactPerson,
		PrimaryContact:   payload.PrimaryContact,
		SecondaryContact: payload.SecondaryContact,
		Email:            payload.Email,
		Address:          payload.Address,
		CommissionRate:   commissionRate,
		BankDetails:      bankDetailsStr,
		Remarks:          payload.Remarks,
		Status:           statusStr,
		CreatedBy:        payload.CreatedBy,
	}

	if err := ctrl.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create introducer " + err.Error()})
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

	var payload struct {
		IntroducerType   string      `json:"introducer_type"`
		Name             string      `json:"name"`
		RegistrationNo   string      `json:"registration_no"`
		ContactPerson    string      `json:"contact_person"`
		PrimaryContact   string      `json:"primary_contact"`
		SecondaryContact string      `json:"secondary_contact"`
		Email            string      `json:"email"`
		Address          string      `json:"address"`
		CommissionRate   interface{} `json:"commission_rate"`
		BankDetails      interface{} `json:"bank_details"`
		Remarks          string      `json:"remarks"`
		Status           interface{} `json:"status"`
		CreatedBy        uint        `json:"created_by"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commissionRate, _ := strconv.ParseFloat(fmt.Sprintf("%v", payload.CommissionRate), 64)

	statusStr := record.Status
	if payload.Status != nil {
		statusStr = fmt.Sprintf("%v", payload.Status)
		if statusStr == "1" {
			statusStr = "Active"
		} else if statusStr == "0" {
			statusStr = "Inactive"
		}
	}

	bankDetailsStr := record.BankDetails
	if payload.BankDetails != nil {
		switch v := payload.BankDetails.(type) {
		case string:
			if v != "" {
				bankDetailsStr = v
			} else {
				bankDetailsStr = "{}"
			}
		default:
			b, _ := json.Marshal(v)
			bankDetailsStr = string(b)
		}
	}
	if bankDetailsStr == "" {
		bankDetailsStr = "{}"
	}

	record.IntroducerType = payload.IntroducerType
	record.Name = payload.Name
	record.RegistrationNo = payload.RegistrationNo
	record.ContactPerson = payload.ContactPerson
	record.PrimaryContact = payload.PrimaryContact
	record.SecondaryContact = payload.SecondaryContact
	record.Email = payload.Email
	record.Address = payload.Address
	record.CommissionRate = commissionRate
	record.BankDetails = bankDetailsStr
	record.Remarks = payload.Remarks
	record.Status = statusStr
	record.CreatedBy = payload.CreatedBy

	if err := ctrl.DB.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update introducer " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

func (ctrl *IntroducerController) Destroy(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&models.Introducer{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete introducer"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Introducer deleted successfully"})
}
