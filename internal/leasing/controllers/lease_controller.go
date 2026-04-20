package controllers

import (
	"fmt"
	"garment-management-backend/internal/leasing/models"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LeaseController struct {
	DB *gorm.DB
}

// CalculateSummary handles POST /api/leasings/calculate
func (ctrl *LeaseController) CalculateSummary(c *gin.Context) {
	var req struct {
		MarketValue   float64 `json:"market_value"`
		LTV           float64 `json:"ltv"`
		InterestRate  float64 `json:"interest_rate"`
		Period        float64 `json:"period"`
		ProductID     uint    `json:"product_id"`
		ProductItemID *uint   `json:"product_item_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	facilityAmount := 0.0
	if req.LTV > 0 {
		facilityAmount = math.Round(req.MarketValue*(req.LTV/100)*100) / 100
	}

	if facilityAmount <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"success":              false,
			"message":              "Insufficient inputs — please check Market Value and LTV.",
			"facility_amount":      "0.00",
			"interest":             "0.00",
			"installment":          "0.00",
			"total_payable":        "0.00",
			"disbursement_charges": "0.00",
			"first_inst_charges":   "0.00",
			"per_inst_charges":     "0.00",
			"net_disbursement":     "0.00",
		})
		return
	}

	var product models.Product
	if err := ctrl.DB.Preload("AdditionalCharges").First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Unit day normalization
	unitDays := map[string]float64{
		"days":       1,
		"per_day":    1,
		"per_days":   1,
		"weeks":      7.5,
		"per_week":   7.5,
		"months":     30,
		"per_month":  30,
		"per_months": 30,
		"year":       360,
		"per_year":   360,
	}

	loanPeriodType := product.LoanPeriodType
	if loanPeriodType == "" {
		loanPeriodType = "months"
	}
	interestPeriodType := product.InterestPeriodType // NOTE: you may need to add InterestPeriodType to Product model
	if interestPeriodType == "" {
		interestPeriodType = "per_month"
	}

	loanUnitDays, ok1 := unitDays[loanPeriodType]
	if !ok1 {
		loanUnitDays = 30
	}

	interestUnitDays, ok2 := unitDays[interestPeriodType]
	if !ok2 {
		interestUnitDays = 30
	}

	totalInterestRate := req.InterestRate * (loanUnitDays / interestUnitDays)

	// Calculate interest
	isReducing := product.InterestMethod == "Reducing Balance" || product.InterestMethod == "reducing_balance"
	var interestAmount, installment float64

	if isReducing && req.Period > 0 {
		periodicRate := totalInterestRate / 100
		n := req.Period
		var emi float64
		if periodicRate > 0 {
			emi = (facilityAmount * periodicRate * math.Pow(1+periodicRate, n)) / (math.Pow(1+periodicRate, n) - 1)
		} else {
			emi = facilityAmount / n
		}
		interestAmount = math.Round(((emi*n)-facilityAmount)*100) / 100
		installment = math.Round(emi*100) / 100
	} else {
		// Flat rate
		interestAmount = math.Round(facilityAmount*(totalInterestRate/100)*req.Period*100) / 100
		if req.Period > 0 {
			installment = math.Round(((facilityAmount + interestAmount) / req.Period) * 100) / 100
		} else {
			installment = 0
		}
	}

	totalPayable := math.Round((facilityAmount+interestAmount)*100) / 100

	// Charges
	disbursementCharges := 0.0
	firstInstCharges := 0.0
	perInstCharges := 0.0

	for _, charge := range product.AdditionalCharges {
		val := charge.Value
		var amt float64
		if strings.ToLower(charge.ValueType) == "percentage" {
			amt = facilityAmount * (val / 100)
		} else {
			amt = val
		}

		dtype := strings.ToLower(charge.DeductionType)
		if strings.Contains(dtype, "disbursement") {
			disbursementCharges += amt
		} else if strings.Contains(dtype, "first") && strings.Contains(dtype, "installment") {
			firstInstCharges += amt
		} else if strings.Contains(dtype, "installment") {
			perInstCharges += amt
		}
	}

	installmentWithCharges := math.Round((installment+perInstCharges)*100) / 100
	netDisbursement := math.Round((facilityAmount-disbursementCharges)*100) / 100

	c.JSON(http.StatusOK, gin.H{
		"success":              true,
		"facility_amount":      fmt.Sprintf("%.2f", facilityAmount),
		"interest":             fmt.Sprintf("%.2f", interestAmount),
		"installment":          fmt.Sprintf("%.2f", installmentWithCharges),
		"total_payable":        fmt.Sprintf("%.2f", totalPayable),
		"disbursement_charges": fmt.Sprintf("%.2f", disbursementCharges),
		"first_inst_charges":   fmt.Sprintf("%.2f", firstInstCharges),
		"per_inst_charges":     fmt.Sprintf("%.2f", perInstCharges),
		"net_disbursement":     fmt.Sprintf("%.2f", netDisbursement),
	})
}
