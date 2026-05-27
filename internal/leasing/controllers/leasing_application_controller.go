package controllers

import (
	"encoding/json"
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

type LeasingApplicationController struct {
	DB *gorm.DB
}

// DraftFormData maps the frontend state to check step completion statuses
type DraftFormData struct {
	CustomerID      string `json:"customer_id"`
	Introducers     []any  `json:"introducers"`
	VehicleMake     string `json:"vehicle_make"`
	VehicleModel    string `json:"vehicle_model"`
	InsuranceCompany string `json:"insurance_company"`
	ProductID       string `json:"product_id"`
	LoanAmount      string `json:"loan_amount"`
	Guarantors      []any  `json:"guarantors"`
	PdcSecurities   []any  `json:"pdc_securities"`
	Cheques         []any  `json:"cheques"`
	OriginalCrNo    string `json:"original_cr_no"`
}

func validateDraft(data []byte) map[int]string {
	var form DraftFormData
	statuses := make(map[int]string)
	for i := 1; i <= 9; i++ {
		statuses[i] = "" // Default pristine (no color)
	}

	if len(data) == 0 {
		return statuses
	}

	json.Unmarshal(data, &form)

	// Step 1: Customer
	if form.CustomerID != "" {
		statuses[1] = "complete" // Green
	} else {
		statuses[1] = "error" // Orange
	}

	// Step 2: Introducers (Optional, but if they visited or added)
	if len(form.Introducers) > 0 {
		statuses[2] = "complete"
	}

	// Step 3: Vehicle Asset
	if form.VehicleMake != "" && form.VehicleModel != "" {
		statuses[3] = "complete"
	} else if form.VehicleMake != "" || form.VehicleModel != "" {
		statuses[3] = "error"
	}

	// Step 4: Insurance
	if form.InsuranceCompany != "" {
		statuses[4] = "complete"
	}

	// Step 5: Lease Details
	if form.ProductID != "" && form.LoanAmount != "" && form.LoanAmount != "0.00" {
		statuses[5] = "complete"
	} else if form.ProductID != "" || form.LoanAmount != "" {
		statuses[5] = "error"
	}

	// Step 6: Guarantors
	if len(form.Guarantors) > 0 {
		statuses[6] = "complete"
	}

	// Step 7: PDC Security
	if len(form.PdcSecurities) > 0 {
		statuses[7] = "complete"
	}

	// Step 8: Cheque Define
	if len(form.Cheques) > 0 {
		statuses[8] = "complete"
	}

	// Step 9: CR & Docs
	if form.OriginalCrNo != "" {
		statuses[9] = "complete"
	}

	return statuses
}

// CreateDraft handles POST /api/v1/leasing-applications/draft
func (ctrl *LeasingApplicationController) CreateDraft(c *gin.Context) {
	var req struct {
		CustomerID          uint            `json:"customer_id" binding:"required"`
		IntroducerID        *uint           `json:"introducer_id"`
		BranchID            *uint           `json:"branch_id"`
		CurrentProgressData json.RawMessage `json:"current_progress_data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app := models.LeasingApplication{
		CustomerID:          req.CustomerID,
		IntroducerID:        req.IntroducerID,
		BranchID:            req.BranchID,
		Status:              "draft",
		CurrentProgressData: string(req.CurrentProgressData),
	}

	if err := ctrl.DB.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create draft"})
		return
	}

	// Generate a temporary Loan Code if needed, or leave for submit phase
	c.JSON(http.StatusCreated, gin.H{
		"message": "Draft created successfully",
		"data":    app,
		"step_statuses": validateDraft(req.CurrentProgressData),
	})
}

// UpdateDraft handles PUT /api/v1/leasing-applications/:id/draft
func (ctrl *LeasingApplicationController) UpdateDraft(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		CurrentProgressData json.RawMessage `json:"current_progress_data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := ctrl.DB.Model(&models.LeasingApplication{}).
		Where("id = ?", id).
		Update("current_progress_data", string(req.CurrentProgressData))
		
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update draft"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Draft not found in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Draft updated successfully",
		"step_statuses": validateDraft(req.CurrentProgressData),
	})
}

// UploadDocument handles POST /api/v1/leasing-applications/:id/upload-document
func (ctrl *LeasingApplicationController) UploadDocument(c *gin.Context) {
	appID := c.Param("id")
	imageType := c.PostForm("image_type")

	// Verify application exists
	var app models.LeasingApplication
	if err := ctrl.DB.First(&app, appID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Leasing application not found"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	safeFilename := fmt.Sprintf("app_%s_%d%s", appID, time.Now().UnixMilli(), ext)
	uploadDir := "./uploads/leasing_documents"

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

	// Create DB Record immediately as per user request
	doc := models.LeasingVehicleDocumentImage{
		LeasingApplicationID: &app.ID,
		ImageType:            imageType,
		ImageURL:             "leasing_documents/" + safeFilename,
	}

	if err := ctrl.DB.Create(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Document uploaded successfully",
		"document_id": doc.ID,
		"url":         doc.ImageURL,
	})
}

// Submit handles POST /api/v1/leasing-applications/:id/submit
// This is a stub that should parse the JSON and populate the real relational tables.
func (ctrl *LeasingApplicationController) Submit(c *gin.Context) {
	id := c.Param("id")

	var fullData struct {
		Vehicle      models.LeasingVehicle      `json:"Vehicle"`
		Loan         models.LeasingLoan         `json:"Loan"`
		Guarantors   []models.LeasingGuarantor  `json:"Guarantors"`
		PdcSecurity  models.PdcSecurity         `json:"PdcSecurity"`
		ChequeDefine models.LeasingChequeDefine `json:"ChequeDefine"`
	}

	if err := c.ShouldBindJSON(&fullData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submit payload: " + err.Error()})
		return
	}
	
	// Start transaction
	tx := ctrl.DB.Begin()

	var app models.LeasingApplication
	if err := tx.First(&app, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	if app.Status != "draft" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application is not in draft state"})
		return
	}

	// Insert transformed data
	fullData.Vehicle.LeasingApplicationID = &app.ID
	if err := tx.Create(&fullData.Vehicle).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save vehicle: " + err.Error()})
		return
	}

	// Link uploaded documents/images to the created vehicle record
	if err := tx.Model(&models.LeasingVehicleDocumentImage{}).
		Where("leasing_application_id = ?", app.ID).
		Update("leasing_vehicle_id", fullData.Vehicle.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link uploaded documents to vehicle: " + err.Error()})
		return
	}

	fullData.Loan.LeasingApplicationID = &app.ID
	fullData.Loan.LeasingVehicleID = &fullData.Vehicle.ID
	if err := tx.Create(&fullData.Loan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save loan: " + err.Error()})
		return
	}

	for _, g := range fullData.Guarantors {
		g.LeasingApplicationID = &app.ID
		if err := tx.Create(&g).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save guarantor: " + err.Error()})
			return
		}
	}

	fullData.PdcSecurity.LeasingApplicationID = &app.ID
	if err := tx.Create(&fullData.PdcSecurity).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save PDC security: " + err.Error()})
		return
	}

	fullData.ChequeDefine.LeasingApplicationID = &app.ID
	if err := tx.Create(&fullData.ChequeDefine).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save Cheque define: " + err.Error()})
		return
	}

	// Update status to pending
	if err := tx.Model(&app).Update("status", "pending").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application status"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Application submitted successfully and status changed to pending",
	})
}

// Get handles GET /api/v1/leasing-applications/:id
func (ctrl *LeasingApplicationController) Get(c *gin.Context) {
	id := c.Param("id")
	var app models.LeasingApplication

	if err := ctrl.DB.
		Preload("Customer").
		Preload("Introducer").
		Preload("Vehicle").
		Preload("Vehicle.Images").
		Preload("Loan").
		Preload("Guarantors").
		Preload("PdcSecurity").
		Preload("PdcSecurity.ChequeDetails").
		Preload("ChequeDefine").
		Preload("ChequeDefine.Items").
		Preload("DocumentImages").
		First(&app, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": app})
}

// GetDrafts handles GET /api/v1/leasing-applications/drafts
func (ctrl *LeasingApplicationController) GetDrafts(c *gin.Context) {
	var drafts []models.LeasingApplication
	query := ctrl.DB.Model(&models.LeasingApplication{}).Where("status = ?", "draft")

	code := c.Query("code")
	nic := c.Query("nic")

	if nic != "" {
		// To search by NIC we need to join the customers table
		query = query.Joins("JOIN customers ON customers.id = leasing_applications.customer_id").
			Where("customers.old_nic ILIKE ? OR customers.new_nic ILIKE ?", "%"+nic+"%", "%"+nic+"%")
	}

	if err := query.Preload("Customer").Find(&drafts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch drafts"})
		return
	}

	// We format the response to match the frontend expectations
	type DraftResponse struct {
		ID                  uint      `json:"ID"`
		DraftCode           string    `json:"draft_code"`
		CustomerID          uint      `json:"customer_id"`
		CurrentProgressData string    `json:"current_progress_data"`
		Status              string    `json:"status"`
		CreatedAt           time.Time `json:"CreatedAt"`
		UpdatedAt           time.Time `json:"UpdatedAt"`
	}

	var response []DraftResponse
	for _, d := range drafts {
		// If code filter is applied, we must check it against the generated draft_code
		draftCode := fmt.Sprintf("LSE-2026-%04d", d.ID)
		
		if code != "" {
			// Simple contains check for the code
			if !strings.Contains(strings.ToLower(draftCode), strings.ToLower(code)) {
				continue
			}
		}

		response = append(response, DraftResponse{
			ID:                  d.ID,
			DraftCode:           draftCode,
			CustomerID:          d.CustomerID,
			CurrentProgressData: d.CurrentProgressData,
			Status:              d.Status,
			CreatedAt:           d.CreatedAt,
			UpdatedAt:           d.UpdatedAt,
		})
	}

	// If response is nil, return empty array instead of null
	if response == nil {
		response = []DraftResponse{}
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// UpdateDraftStep handles PUT /api/v1/leasing-applications/:id/draft/step/:step_name
func (ctrl *LeasingApplicationController) UpdateDraftStep(c *gin.Context) {
	id := c.Param("id")
	stepName := c.Param("step_name")

	// Read request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Parse step data into a map
	var stepData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &stepData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
		return
	}

	// Load existing application
	var app models.LeasingApplication
	if err := ctrl.DB.First(&app, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Draft not found"})
		return
	}

	// Parse existing progress data
	var existingData map[string]interface{}
	if app.CurrentProgressData != "" {
		if err := json.Unmarshal([]byte(app.CurrentProgressData), &existingData); err != nil {
			existingData = make(map[string]interface{})
		}
	} else {
		existingData = make(map[string]interface{})
	}

	// Merge stepData into existingData
	for k, v := range stepData {
		existingData[k] = v
	}

	// Marshal back to JSON
	updatedBytes, err := json.Marshal(existingData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize progress data"})
		return
	}

	// Prepare fields for updates
	updateFields := map[string]interface{}{
		"current_progress_data": string(updatedBytes),
	}

	// Sync customer_id and introducer_id if present
	if cidVal, ok := stepData["customer_id"]; ok {
		if cidNum, err := parseUint(cidVal); err == nil {
			updateFields["customer_id"] = cidNum
		}
	}
	if iidVal, ok := stepData["introducer_id"]; ok {
		if iidNum, err := parseUint(iidVal); err == nil {
			updateFields["introducer_id"] = iidNum
		} else if iidVal == nil || iidVal == "" {
			updateFields["introducer_id"] = nil
		}
	}

	if err := ctrl.DB.Model(&app).Updates(updateFields).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save draft step"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       fmt.Sprintf("Step '%s' draft updated successfully", stepName),
		"step_statuses": validateDraft(updatedBytes),
	})
}

func parseUint(val interface{}) (uint, error) {
	switch v := val.(type) {
	case float64:
		return uint(v), nil
	case int:
		return uint(v), nil
	case string:
		var n uint
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
			return n, nil
		}
	}
	return 0, fmt.Errorf("invalid type for uint")
}
