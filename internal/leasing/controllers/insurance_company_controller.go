package controllers

import (
	"fmt"
	"garment-management-backend/internal/leasing/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InsuranceCompanyController struct {
	DB *gorm.DB
}

func (ctrl *InsuranceCompanyController) Index(c *gin.Context) {
	var records []models.InsuranceCompany
	if err := ctrl.DB.Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch insurance companies"})
		return
	}
	c.JSON(http.StatusOK, records)
}

func (ctrl *InsuranceCompanyController) Store(c *gin.Context) {
	var payload struct {
		CompanyCode          string      `json:"company_code"`
		CompanyName          string      `json:"company_name"`
		HeadOfficeAddress    string      `json:"head_office_address"`
		ContactPerson        string      `json:"contact_person"`
		ContactMobile        string      `json:"contact_mobile"`
		ContactEmail         string      `json:"contact_email"`
		ContactPerson2       string      `json:"contact_person2"`
		ContactPerson2Mobile string      `json:"contact_person2_mobile"`
		ContactPerson2Email  string      `json:"contact_person2_email"`
		CommisionRate        interface{} `json:"commision_rate"`
		BankAccountNo        string      `json:"bank_account_no"`
		BankAccountName      string      `json:"bank_account_name"`
		BankName             string      `json:"bank_name"`
		Status               interface{} `json:"status"`
		CreatedBy            uint        `json:"created_by"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commissionRate, _ := strconv.ParseFloat(fmt.Sprintf("%v", payload.CommisionRate), 64)

	statusStr := "Active"
	if payload.Status != nil {
		statusStr = fmt.Sprintf("%v", payload.Status)
		if statusStr == "1" || statusStr == "true" {
			statusStr = "Active"
		} else if statusStr == "0" || statusStr == "false" {
			statusStr = "Inactive"
		}
	}

	req := models.InsuranceCompany{
		CompanyCode:          payload.CompanyCode,
		CompanyName:          payload.CompanyName,
		HeadOfficeAddress:    payload.HeadOfficeAddress,
		ContactPerson:        payload.ContactPerson,
		ContactMobile:        payload.ContactMobile,
		ContactEmail:         payload.ContactEmail,
		ContactPerson2:       payload.ContactPerson2,
		ContactPerson2Mobile: payload.ContactPerson2Mobile,
		ContactPerson2Email:  payload.ContactPerson2Email,
		CommisionRate:        commissionRate,
		BankAccountNo:        payload.BankAccountNo,
		BankAccountName:      payload.BankAccountName,
		BankName:             payload.BankName,
		Status:               statusStr,
		CreatedBy:            payload.CreatedBy,
	}

	if err := ctrl.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create insurance company"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func (ctrl *InsuranceCompanyController) Update(c *gin.Context) {
	id := c.Param("id")
	var record models.InsuranceCompany
	if err := ctrl.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Insurance company not found"})
		return
	}

	var payload struct {
		CompanyCode          string      `json:"company_code"`
		CompanyName          string      `json:"company_name"`
		HeadOfficeAddress    string      `json:"head_office_address"`
		ContactPerson        string      `json:"contact_person"`
		ContactMobile        string      `json:"contact_mobile"`
		ContactEmail         string      `json:"contact_email"`
		ContactPerson2       string      `json:"contact_person2"`
		ContactPerson2Mobile string      `json:"contact_person2_mobile"`
		ContactPerson2Email  string      `json:"contact_person2_email"`
		CommisionRate        interface{} `json:"commision_rate"`
		BankAccountNo        string      `json:"bank_account_no"`
		BankAccountName      string      `json:"bank_account_name"`
		BankName             string      `json:"bank_name"`
		Status               interface{} `json:"status"`
		CreatedBy            uint        `json:"created_by"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commissionRate, _ := strconv.ParseFloat(fmt.Sprintf("%v", payload.CommisionRate), 64)

	statusStr := record.Status
	if payload.Status != nil {
		statusStr = fmt.Sprintf("%v", payload.Status)
		if statusStr == "1" || statusStr == "true" {
			statusStr = "Active"
		} else if statusStr == "0" || statusStr == "false" {
			statusStr = "Inactive"
		}
	}

	record.CompanyCode = payload.CompanyCode
	record.CompanyName = payload.CompanyName
	record.HeadOfficeAddress = payload.HeadOfficeAddress
	record.ContactPerson = payload.ContactPerson
	record.ContactMobile = payload.ContactMobile
	record.ContactEmail = payload.ContactEmail
	record.ContactPerson2 = payload.ContactPerson2
	record.ContactPerson2Mobile = payload.ContactPerson2Mobile
	record.ContactPerson2Email = payload.ContactPerson2Email
	record.CommisionRate = commissionRate
	record.BankAccountNo = payload.BankAccountNo
	record.BankAccountName = payload.BankAccountName
	record.BankName = payload.BankName
	record.Status = statusStr
	record.CreatedBy = payload.CreatedBy

	if err := ctrl.DB.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update insurance company"})
		return
	}

	c.JSON(http.StatusOK, record)
}

func (ctrl *InsuranceCompanyController) Destroy(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&models.InsuranceCompany{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete insurance company"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Insurance company deleted successfully"})
}
