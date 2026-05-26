package controllers

import (
	"fmt"
	"garment-management-backend/internal/leasing/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CustomerController struct {
	DB *gorm.DB
}

// resolveNics converts between old (10-char + V) and new (12-digit) NIC formats.
func resolveNics(inputNic string) (newNic, oldNic string) {
	nic := strings.ToUpper(strings.TrimSpace(inputNic))
	switch len(nic) {
	case 10: // Old format e.g. 853400937V
		oldNic = nic
		base := nic[:9]
		newNic = "19" + base[:5] + "0" + base[5:9]
	case 12: // New format e.g. 198534000937
		newNic = nic
		if strings.HasPrefix(nic, "19") {
			oldBase := nic[2:7] + nic[8:12]
			oldNic = oldBase + "V"
		}
	default:
		newNic = nic
	}
	return
}

// Store handles POST /api/customers
func (ctrl *CustomerController) Store(c *gin.Context) {
	var req struct {
		CustomerID         string `json:"customer_id"`
		Title              string `json:"title"`
		FullName           string `json:"full_name"`
		FirstName          string `json:"first_name"`
		LastName           string `json:"last_name"`
		NameWithInitials   string `json:"name_with_initials"`
		NIC                string `json:"nic"`
		DOB                string `json:"dob"`
		Gender             string `json:"gender"`
		Status             string `json:"status"`
		PerAddressLine1    string `json:"permanent_address_line1"`
		PerAddressLine2    string `json:"permanent_address_line2"`
		PerAddressLine3    string `json:"permanent_address_line3"`
		PostalAddressLine1 string `json:"postal_address_line1"`
		PostalAddressLine2 string `json:"postal_address_line2"`
		PostalAddressLine3 string `json:"postal_address_line3"`
		Province           string `json:"province"`
		City               string `json:"city"`
		MobilePrimary      string `json:"mobile_primary"`
		MobileSecondary    string `json:"mobile_secondary"`
		Landline           string `json:"landline"`
		Email              string `json:"email"`
		Remarks            string `json:"remarks"`

		Occupations []struct {
			EngagementType     string `json:"engagementType"`
			Position           string `json:"position"`
			BusinessName       string `json:"businessName"`
			RegistrationNumber string `json:"registrationNumber"`
			NatureOfBusiness   string `json:"natureOfBusiness"`
			NetMonthlyIncome   string `json:"netMonthlyIncome"`
			StartDate          string `json:"startDate"`
			EndDate            string `json:"endDate"`
		} `json:"occupations"`

		BankAccounts []struct {
			Bank          string `json:"bank"`
			AccountNumber string `json:"accountNumber"`
			Beneficiary   string `json:"beneficiary"`
			Type          string `json:"type"`
			Branch        string `json:"branch"`
		} `json:"bank_accounts"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Resolve NIC formats
	newNic, oldNic := resolveNics(req.NIC)

	// Begin transaction across multiple tables
	tx := ctrl.DB.Begin()

	customer := models.Customer{
		CustomerCode:       req.CustomerID,
		Title:              req.Title,
		FullName:           req.FullName,
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		NameWithInitials:   req.NameWithInitials,
		NewNic:             newNic,
		OldNic:             oldNic,
		DOB:                req.DOB,
		Gender:             req.Gender,
		Status:             req.Status,
		ContactNo:          req.MobilePrimary,
		ContactNo2:         req.MobileSecondary,
		Landline:           req.Landline,
		Email:              req.Email,
		Remarks:            req.Remarks,
		Province:           req.Province,
		City:               req.City,
		PerAddressLine1:    req.PerAddressLine1,
		PerAddressLine2:    req.PerAddressLine2,
		PerAddressLine3:    req.PerAddressLine3,
		PostalAddressLine1: req.PostalAddressLine1,
		PostalAddressLine2: req.PostalAddressLine2,
		PostalAddressLine3: req.PostalAddressLine3,
	}

	if err := tx.Create(&customer).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer" + err.Error()})
		return
	}

	// Save occupations
	for _, occ := range req.Occupations {
		o := models.CustomerOccupation{
			CustomerID:       customer.ID,
			Type:             occ.EngagementType,
			Designation:      occ.Position,
			BRNo:             occ.RegistrationNumber,
			BusinessName:     occ.BusinessName,
			NatureOfBusiness: occ.NatureOfBusiness,
			FromDate:         occ.StartDate,
			ToDate:           occ.EndDate,
		}
		if err := tx.Create(&o).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save occupation"})
			return
		}
	}

	// Save bank accounts
	for _, ba := range req.BankAccounts {
		b := models.CustomerBankAccount{
			CustomerID:    customer.ID,
			BankName:      ba.Bank,
			AccountName:   ba.Beneficiary,
			AccountNumber: ba.AccountNumber,
			Branch:        ba.Branch,
			AccountType:   ba.Type,
		}
		if err := tx.Create(&b).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save bank account"})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Customer created successfully",
		"customer_id": customer.CustomerCode,
	})
}

// Index handles GET /api/customers?code=...&nic=...&status=...&query=...
func (ctrl *CustomerController) Index(c *gin.Context) {
	var customers []models.Customer
	db := ctrl.DB

	// Filters
	if code := c.Query("code"); code != "" {
		db = db.Where("customer_code ILIKE ?", "%"+code+"%")
	}
	if nic := c.Query("nic"); nic != "" {
		db = db.Where("new_nic ILIKE ? OR old_nic ILIKE ?", "%"+nic+"%", "%"+nic+"%")
	}
	if status := c.Query("status"); status != "" {
		db = db.Where("status = ?", status)
	}
	if query := c.Query("query"); query != "" {
		likeQuery := "%" + query + "%"
		db = db.Where("full_name ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ? OR customer_code ILIKE ? OR new_nic ILIKE ? OR old_nic ILIKE ?",
			likeQuery, likeQuery, likeQuery, likeQuery, likeQuery, likeQuery)
	}

	if err := db.Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customers})
}

// Get handles GET /api/customers/:id
func (ctrl *CustomerController) Get(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer

	if err := ctrl.DB.Preload("Occupations").Preload("BankAccounts").Preload("Documents").First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

// UpdateStatus handles POST /api/customers/:id/status
func (ctrl *CustomerController) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Model(&models.Customer{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

// UpdateLocation handles POST /api/customers/:id/location
func (ctrl *CustomerController) UpdateLocation(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Model(&models.Customer{}).Where("id = ?", id).Updates(map[string]interface{}{
		"latitude":  req.Latitude,
		"longitude": req.Longitude,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Location updated successfully"})
}

// GenerateID handles GET /api/customers/generate-id
// Returns the next sequential customer ID in the format CUS-XXXXX
func (ctrl *CustomerController) GenerateID(c *gin.Context) {
	var maxID int64
	ctrl.DB.Model(&models.Customer{}).Count(&maxID)
	next := maxID + 1
	code := fmt.Sprintf("CUS-%05d", next)
	c.JSON(http.StatusOK, gin.H{"customer_id": code})
}

// GetCities handles GET /api/locations/cities?province=Western
func (ctrl *CustomerController) GetCities(c *gin.Context) {
	province := c.Query("province")

	cityMap := map[string][]string{
		"Western":       {"Colombo", "Gampaha", "Kalutara"},
		"Central":       {"Kandy", "Matale", "Nuwara Eliya"},
		"Southern":      {"Galle", "Matara", "Hambantota"},
		"Northern":      {"Jaffna", "Kilinochchi", "Mannar", "Mullaitivu", "Vavuniya"},
		"Eastern":       {"Ampara", "Batticaloa", "Trincomalee"},
		"North Western": {"Kurunegala", "Puttalam"},
		"North Central": {"Anuradhapura", "Polonnaruwa"},
		"Uva":           {"Badulla", "Monaragala"},
		"Sabaragamuwa":  {"Kegalle", "Ratnapura"},
	}

	cities, ok := cityMap[province]
	if !ok || province == "" {
		// Return all if province not found
		all := []string{}
		for _, v := range cityMap {
			all = append(all, v...)
		}
		c.JSON(http.StatusOK, gin.H{"cities": all})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cities": cities})
}

// UploadDocument handles POST /api/customers/:id/documents
// Accepts multipart/form-data with fields: category, file
func (ctrl *CustomerController) UploadDocument(c *gin.Context) {
	customerID := c.Param("id")
	category := c.PostForm("category")

	// Verify customer exists first
	var customer models.Customer
	if err := ctrl.DB.First(&customer, customerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Build a unique filename: customerID_timestamp_originalname
	ext := filepath.Ext(header.Filename)
	safeFilename := fmt.Sprintf("%s_%d%s", customerID, time.Now().UnixMilli(), ext)
	uploadDir := "./uploads/customers"

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create upload directory"})
		return
	}

	savePath := filepath.Join(uploadDir, safeFilename)
	out, err := os.Create(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
		return
	}

	// Store metadata in DB
	doc := models.CustomerDocument{
		CustomerID:  customer.ID,
		Description: category,
		Path:        "customers/" + safeFilename, // Relative path under /uploads
	}

	if err := ctrl.DB.Create(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Document uploaded",
		"document_id": doc.ID,
	})
}

// Search handles GET /api/customers/search?query=val
func (ctrl *CustomerController) Search(c *gin.Context) {
	query := c.Query("query")
	var customers []models.Customer

	if query == "" {
		// Return empty list if no query
		c.JSON(http.StatusOK, gin.H{"data": []models.Customer{}})
		return
	}

	likeQuery := "%" + query + "%"
	// Search by NewNic, OldNic, FirstName, LastName, FullName, CustomerCode
	if err := ctrl.DB.Where("new_nic ILIKE ? OR old_nic ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ? OR full_name ILIKE ? OR customer_code ILIKE ?", 
		likeQuery, likeQuery, likeQuery, likeQuery, likeQuery, likeQuery).
		Limit(20).
		Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search customers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customers})
}

// GetBankAccounts handles GET /api/customers/:id/bank-accounts
func (ctrl *CustomerController) GetBankAccounts(c *gin.Context) {
	id := c.Param("id")
	var accounts []models.CustomerBankAccount

	if err := ctrl.DB.Where("customer_id = ?", id).Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bank accounts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": accounts})
}
