package controllers

import (
	"garment-management-backend/internal/leasing/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct {
	DB *gorm.DB
}

// Store handles the POST /v1/leasing/products endpoint
func (ctrl *ProductController) Store(c *gin.Context) {
	// Let's bind the incoming payload strictly or mapping to our JSON format
	// given by the frontend.
	var req struct {
		Name               string `json:"name"`
		Code               string `json:"code"`
		InterestMethodID   int    `json:"interest_method_id"` // could map to interest method or stored directly
		Periods            int    `json:"periods"`
		GuarantorsRequired int    `json:"guarantors_required"`

		Configurations []struct {
			Label        string  `json:"label"` // ProductItemName
			MinLoan      float64 `json:"minLoan"`
			MaxLoan      float64 `json:"maxLoan"`
			MinInt       float64 `json:"minInt"`
			MaxInt       float64 `json:"maxInt"`
			MinPeriod    int     `json:"minPeriod"`
			MaxPeriod    int     `json:"maxPeriod"`
			Guarantors   int     `json:"guarantors"`
			PenaltyType  string  `json:"penaltyType"`
			PenaltyRate  float64 `json:"penaltyRate"`
		} `json:"configurations"`

		Charges []struct {
			Description string  `json:"description"`
			Amount      float64 `json:"amount"`
			Type        string  `json:"type"`      // fixed | percentage
			Deduction   string  `json:"deduction"` // on_loan_disbursement | as_first_installment
		} `json:"charges"`

		Documents []struct {
			Name   string `json:"name"`
			Status string `json:"status"` // Required | Optional
		} `json:"documents"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Begin Transaction since we are saving across multiple tables
	tx := ctrl.DB.Begin()

	// Map generic request to specific gorm models
	product := models.Product{
		ProductName:    req.Name,
		ProductCode:    req.Code,
		GuaranteeCount: req.GuarantorsRequired,
		// Example mapping (some fields maybe omitted or defaulted depending on actual logic)
		// For now we map what we have in the payload
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	// Map configurations
	for _, conf := range req.Configurations {
		item := models.ProductHasItem{
			ProductID:               product.ID,
			ProductItemName:         conf.Label,
			MinimumLoanAmount:       conf.MinLoan,
			MaximumLoanAmount:       conf.MaxLoan,
			MinimumInterest:         conf.MinInt,
			MaximumInterest:         conf.MaxInt,
			MinimumLoanPeriod:       conf.MinPeriod,
			MaximumLoanPeriod:       conf.MaxPeriod,
			RequiredGuaranteeCount:  conf.Guarantors,
			PenaltyApplyType:        conf.PenaltyType,
			PenaltyPercentage:       conf.PenaltyRate,
		}
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create configurations"})
			return
		}
	}

	// Map charges
	for _, charge := range req.Charges {
		chargeModel := models.ProductAdditionalCharges{
			ProductID:     product.ID,
			Description:   charge.Description,
			Value:         charge.Amount,
			ValueType:     charge.Type,
			DeductionType: charge.Deduction,
		}
		if err := tx.Create(&chargeModel).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create charges"})
			return
		}
	}

	// Map documents
	for _, doc := range req.Documents {
		d := models.ProductRequiredDocuments{
			ProductID:      product.ID,
			Name:           doc.Name,
			RequiredStatus: doc.Status,
		}
		if err := tx.Create(&d).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create required documents"})
			return
		}
	}

	// Commit Transaction
	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{"message": "Leasing product created successfully", "product_id": product.ID})
}

// Get handles GET /v1/leasing/products/:id
func (ctrl *ProductController) Get(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := ctrl.DB.Preload("ProductHasItems").
		Preload("AdditionalCharges").
		Preload("RequiredDocuments").
		First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateStatus handles POST /v1/leasing/products/status
func (ctrl *ProductController) UpdateStatus(c *gin.Context) {
	var req struct {
		ID     int    `json:"id"`
		Status string `json:"status"` // e.g., "Active" or "Inactive", frontend uses boolean or string depends, typical is string or bool. If we assume "status", string is safer.
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Model(&models.Product{}).Where("id = ?", req.ID).Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

// Index handles GET /v1/leasing/products
func (ctrl *ProductController) Index(c *gin.Context) {
	var products []models.Product

	if err := ctrl.DB.Preload("ProductHasItems").
		Preload("AdditionalCharges").
		Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	type ProductListItem struct {
		ID             uint   `json:"id"`
		ProductName    string `json:"product_name"`
		ProductCode    string `json:"product_code"`
		InterestMethod string `json:"interest_method"`
		LoanPeriodType string `json:"loan_period_type"`
		Status         string `json:"status"`
		ItemsCount     int    `json:"items_count"`
		ChargesCount   int    `json:"charges_count"`
	}

	var response []ProductListItem = make([]ProductListItem, 0)
	for _, p := range products {
		response = append(response, ProductListItem{
			ID:             p.ID,
			ProductName:    p.ProductName,
			ProductCode:    p.ProductCode,
			InterestMethod: p.InterestMethod,
			LoanPeriodType: p.LoanPeriodType,
			Status:         p.Status,
			ItemsCount:     len(p.ProductHasItems),
			ChargesCount:   len(p.AdditionalCharges),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
