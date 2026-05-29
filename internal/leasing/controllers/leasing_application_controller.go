package controllers

import (
	"encoding/json"
	"fmt"
	"garment-management-backend/internal/leasing/models"
	"garment-management-backend/internal/leasing/validation"
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

// isStepTouched returns true if the user has entered some data for a specific step.
func isStepTouched(stepNum int, fields map[string]interface{}) bool {
	switch stepNum {
	case 1:
		cid, ok := fields["customer_id"]
		if ok && cid != nil && fmt.Sprintf("%v", cid) != "" && fmt.Sprintf("%v", cid) != "0" {
			return true
		}
	case 2:
		intros, ok := fields["introducers"].([]interface{})
		if ok && len(intros) > 0 {
			return true
		}
	case 3:
		makeID, _ := fields["vehicle_make_id"]
		modelID, _ := fields["vehicle_model_id"]
		chassisNo, _ := fields["chassis_no"]
		regNo, _ := fields["reg_no"]
		if (makeID != nil && fmt.Sprintf("%v", makeID) != "" && fmt.Sprintf("%v", makeID) != "0") ||
			(modelID != nil && fmt.Sprintf("%v", modelID) != "" && fmt.Sprintf("%v", modelID) != "0") ||
			(chassisNo != nil && fmt.Sprintf("%v", chassisNo) != "") ||
			(regNo != nil && fmt.Sprintf("%v", regNo) != "") {
			return true
		}
	case 4:
		company, _ := fields["insurance_company"]
		amount, _ := fields["insurance_amount"]
		premium, _ := fields["insurance_premium"]
		startDate, _ := fields["insurance_start_date"]
		endDate, _ := fields["insurance_expiry_date"]
		if (company != nil && fmt.Sprintf("%v", company) != "") ||
			(amount != nil && fmt.Sprintf("%v", amount) != "" && fmt.Sprintf("%v", amount) != "0" && fmt.Sprintf("%v", amount) != "0.00") ||
			(premium != nil && fmt.Sprintf("%v", premium) != "" && fmt.Sprintf("%v", premium) != "0" && fmt.Sprintf("%v", premium) != "0.00") ||
			(startDate != nil && fmt.Sprintf("%v", startDate) != "") ||
			(endDate != nil && fmt.Sprintf("%v", endDate) != "") {
			return true
		}
	case 5:
		prodID, _ := fields["product_id"]
		loanAmt, _ := fields["loan_amount"]
		if (prodID != nil && fmt.Sprintf("%v", prodID) != "" && fmt.Sprintf("%v", prodID) != "0") ||
			(loanAmt != nil && fmt.Sprintf("%v", loanAmt) != "" && fmt.Sprintf("%v", loanAmt) != "0" && fmt.Sprintf("%v", loanAmt) != "0.00") {
			return true
		}
	case 6:
		guars, ok := fields["guarantors"].([]interface{})
		if ok && len(guars) > 0 {
			return true
		}
	case 7:
		secType, _ := fields["pdc_security_type"].(string)
		refDetails, _ := fields["pdc_reference_details"].(string)
		bankID, _ := fields["pdc_bank_id"]
		chequeNo, _ := fields["pdc_cheque_no"].(string)
		bookDate, _ := fields["pdc_book_date"].(string)
		if refDetails != "" ||
			secType == "Cheque" ||
			secType == "CR Book" ||
			(bankID != nil && fmt.Sprintf("%v", bankID) != "" && fmt.Sprintf("%v", bankID) != "0") ||
			chequeNo != "" ||
			bookDate != "" {
			return true
		}
	case 8:
		cheques, ok := fields["cheques"].([]interface{})
		if ok && len(cheques) > 0 {
			return true
		}
	case 9:
		crNo, _ := fields["original_cr_no"].(string)
		serialNo, _ := fields["cr_serial_no"].(string)
		urlFront, _ := fields["url_cr_front"].(string)
		urlBack, _ := fields["url_cr_back"].(string)
		urlInvoice, _ := fields["url_invoice"].(string)
		if crNo != "" || serialNo != "" || urlFront != "" || urlBack != "" || urlInvoice != "" {
			return true
		}
	}
	return false
}

// validateDraft evaluates the completion status for each leasing application wizard step using strategy validation
func (ctrl *LeasingApplicationController) validateDraft(data []byte) map[int]string {
	statuses := make(map[int]string)
	for i := 1; i <= 9; i++ {
		statuses[i] = "" // Default pristine (no color)
	}

	if len(data) == 0 {
		return statuses
	}

	var progressData map[string]interface{}
	if err := json.Unmarshal(data, &progressData); err != nil {
		return statuses
	}

	stepNames := map[int]string{
		1: "step-customer",
		2: "step-introducers",
		3: "step-vehicle",
		4: "step-insurance",
		5: "step-lease-details",
		6: "step-guarantors",
		7: "step-pdc-security",
		8: "step-cheque-define",
		9: "step-documents",
	}

	for stepNum, sName := range stepNames {
		stepFields := getStepFieldsMap(stepNum, progressData)
		if !isStepTouched(stepNum, stepFields) {
			statuses[stepNum] = ""
			continue
		}

		errs := validation.ValidateStep(ctrl.DB, sName, stepFields)
		if len(errs) > 0 {
			statuses[stepNum] = "error"
		} else {
			statuses[stepNum] = "complete"
		}
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
		"step_statuses": ctrl.validateDraft(req.CurrentProgressData),
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
		"step_statuses": ctrl.validateDraft(req.CurrentProgressData),
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

	// Clean up any existing PDC Security and detail records for this leasing application to avoid conflicts
	var existingPdc models.PdcSecurity
	if err := tx.Where("leasing_application_id = ?", app.ID).First(&existingPdc).Error; err == nil {
		tx.Unscoped().Where("pdc_security_id = ?", existingPdc.ID).Delete(&models.PdcChequeDetail{})
		tx.Unscoped().Where("pdc_security_id = ?", existingPdc.ID).Delete(&models.PdcCrBookDetail{})
		tx.Unscoped().Where("pdc_security_id = ?", existingPdc.ID).Delete(&models.PdcDeedDetail{})
		tx.Unscoped().Delete(&existingPdc)
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
		Preload("Vehicle.Color").
		Preload("Vehicle.Images").
		Preload("Loan").
		Preload("Guarantors").
		Preload("PdcSecurity").
		Preload("PdcSecurity.ChequeDetails").
		Preload("PdcSecurity.ChequeDetails.Bank").
		Preload("PdcSecurity.CrBookDetails").
		Preload("PdcSecurity.DeedDetails").
		Preload("ChequeDefine").
		Preload("ChequeDefine.Items").
		Preload("DocumentImages").
		First(&app, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// Parse existing progress data to validate each step
	var progressData map[string]interface{}
	if app.CurrentProgressData != "" {
		json.Unmarshal([]byte(app.CurrentProgressData), &progressData)
	}

	// Step names mapping for verification
	stepNames := map[int]string{
		1: "step-customer",
		2: "step-introducers",
		3: "step-vehicle",
		4: "step-insurance",
		5: "step-lease-details",
		6: "step-guarantors",
		7: "step-pdc-security",
		8: "step-cheque-define",
		9: "step-documents",
	}

	stepErrors := make(map[string]map[string]string)
	if progressData != nil {
		for stepNum, sName := range stepNames {
			stepFields := getStepFieldsMap(stepNum, progressData)
			errs := validation.ValidateStep(ctrl.DB, sName, stepFields)
			if len(errs) > 0 {
				stepErrors[sName] = errs
			} else {
				stepErrors[sName] = make(map[string]string)
			}
		}
	} else {
		// Initialize empty errors map for all steps
		for _, sName := range stepNames {
			stepErrors[sName] = make(map[string]string)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":          app,
		"step_statuses": ctrl.validateDraft([]byte(app.CurrentProgressData)),
		"step_errors":   stepErrors,
	})
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

	// Run Strategy validation
	stepErrors := validation.ValidateStep(ctrl.DB, stepName, stepData)

	c.JSON(http.StatusOK, gin.H{
		"message":       fmt.Sprintf("Step '%s' draft updated successfully", stepName),
		"errors":        stepErrors,
		"step_statuses": ctrl.validateDraft(updatedBytes),
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

// GetPdcSecurity handles GET /api/v1/leasing-applications/:id/pdc-security
func (ctrl *LeasingApplicationController) GetPdcSecurity(c *gin.Context) {
	id := c.Param("id")
	var pdc models.PdcSecurity
	if err := ctrl.DB.
		Preload("ChequeDetails").
		Preload("ChequeDetails.Bank").
		Preload("CrBookDetails").
		Preload("DeedDetails").
		Where("leasing_application_id = ?", id).
		First(&pdc).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "PDC Security not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch PDC security: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": pdc})
}

// UpdatePdcSecurity handles PUT /api/v1/leasing-applications/:id/pdc-security
func (ctrl *LeasingApplicationController) UpdatePdcSecurity(c *gin.Context) {
	id := c.Param("id")
	
	// Parse input payload
	var req models.PdcSecurity
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Start transaction
	tx := ctrl.DB.Begin()

	// Verify application exists
	var app models.LeasingApplication
	if err := tx.First(&app, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Leasing application not found"})
		return
	}

	// Check if a PdcSecurity already exists
	var existingPdc models.PdcSecurity
	hasExisting := true
	if err := tx.Where("leasing_application_id = ?", id).First(&existingPdc).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			hasExisting = false
		} else {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search existing PDC security: " + err.Error()})
			return
		}
	}

	var pdcID uint
	if hasExisting {
		// Delete existing details records to prevent duplicates / orphans
		tx.Unscoped().Where("pdc_security_id = ?", existingPdc.ID).Delete(&models.PdcChequeDetail{})
		tx.Unscoped().Where("pdc_security_id = ?", existingPdc.ID).Delete(&models.PdcCrBookDetail{})
		tx.Unscoped().Where("pdc_security_id = ?", existingPdc.ID).Delete(&models.PdcDeedDetail{})

		// Update parent
		existingPdc.PdcSecurityType = req.PdcSecurityType
		if err := tx.Save(&existingPdc).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update PDC security: " + err.Error()})
			return
		}
		pdcID = existingPdc.ID
	} else {
		// Create new PdcSecurity
		req.LeasingApplicationID = &app.ID
		// Clear nested details so GORM doesn't try to automatically insert them with potential issues
		chequeDetails := req.ChequeDetails
		crBookDetails := req.CrBookDetails
		deedDetails := req.DeedDetails

		req.ChequeDetails = nil
		req.CrBookDetails = nil
		req.DeedDetails = nil

		if err := tx.Create(&req).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create PDC security: " + err.Error()})
			return
		}
		pdcID = req.ID

		req.ChequeDetails = chequeDetails
		req.CrBookDetails = crBookDetails
		req.DeedDetails = deedDetails
	}

	// Insert details based on the selected security type
	if req.PdcSecurityType == "Cheque" {
		for i := range req.ChequeDetails {
			req.ChequeDetails[i].PdcSecurityID = &pdcID
			if err := tx.Create(&req.ChequeDetails[i]).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cheque details: " + err.Error()})
				return
			}
		}
	} else if req.PdcSecurityType == "CR Book" || req.PdcSecurityType == "CR Book (Certificate of Registration)" {
		for i := range req.CrBookDetails {
			req.CrBookDetails[i].PdcSecurityID = &pdcID
			if err := tx.Create(&req.CrBookDetails[i]).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save CR Book details: " + err.Error()})
				return
			}
		}
	} else if req.PdcSecurityType == "Deed" || req.PdcSecurityType == "Deed (Signed Contract)" {
		for i := range req.DeedDetails {
			req.DeedDetails[i].PdcSecurityID = &pdcID
			if err := tx.Create(&req.DeedDetails[i]).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save Deed details: " + err.Error()})
				return
			}
		}
	}

	tx.Commit()

	// Fetch full saved object to return to caller
	var updatedPdc models.PdcSecurity
	ctrl.DB.
		Preload("ChequeDetails").
		Preload("ChequeDetails.Bank").
		Preload("CrBookDetails").
		Preload("DeedDetails").
		First(&updatedPdc, pdcID)

	c.JSON(http.StatusOK, gin.H{
		"message": "PDC Security details saved successfully",
		"data":    updatedPdc,
	})
}

func getStepFieldsMap(step int, data map[string]interface{}) map[string]interface{} {
	fields := make(map[string]interface{})

	copyKeys := func(keys []string) {
		for _, k := range keys {
			if v, ok := data[k]; ok {
				fields[k] = v
			}
		}
	}

	switch step {
	case 1:
		if v, ok := data["customer_id"]; ok {
			fields["customer_id"] = v
		} else if v, ok := data["customer_db_id"]; ok {
			fields["customer_id"] = v
		}
		copyKeys([]string{"customer_code", "customer_name", "bank_account_id"})
	case 2:
		copyKeys([]string{"introducers"})
	case 3:
		copyKeys([]string{
			"vehicle_type", "vehicle_type_id", "vehicle_make", "vehicle_make_id",
			"vehicle_model", "vehicle_model_id", "vehicle_status", "engine_cc",
			"chassis_no", "manu_year", "color", "color_id", "usage_type", "manu_country",
			"body_type", "equipment", "reg_year", "reg_no", "registration_no", "valuation_no",
			"market_value", "forced_value", "invoice_value", "supplier_name",
			"supplier_address", "supplier_mobile", "supplier_id", "supplier_rno",
		})
	case 4:
		copyKeys([]string{
			"insurance_company", "insurance_amount", "insurance_premium",
			"insurance_start_date", "insurance_expiry_date",
		})
	case 5:
		copyKeys([]string{
			"product_id", "product_item", "product_item_id", "marketing_executive_id",
			"inspection_date", "loan_amount", "period", "interest_rate",
			"installments_total", "total_interest", "total_payable",
			"tcc_collection_date", "bank_id", "branch_id", "account_number",
			"bank_account_id", "ltv", "disburse_amount", "installment_amount",
			"other_charges_total", "other_charges_on_disburse",
			"other_charges_on_first_installment", "other_charges_on_every_installments",
			"required_guarantor_count",
		})
	case 6:
		copyKeys([]string{"guarantors", "required_guarantor_count"})
	case 7:
		copyKeys([]string{
			"pdc_security_type", "pdc_cheque_status", "pdc_bank_id",
			"pdc_cheque_date", "pdc_cheque_no", "pdc_ownership",
			"pdc_book_date", "pdc_reference_details",
		})
	case 8:
		copyKeys([]string{"cheques"})
	case 9:
		copyKeys([]string{
			"original_cr_no", "duplicate_key", "documents", "cr_serial_no",
			"url_cr_front", "url_cr_back", "url_invoice", "url_valuation",
		})
	}
	return fields
}
