package models

import "gorm.io/gorm"

// ProductHasItem maps to the `product_has_items` table
type ProductHasItem struct {
	gorm.Model // Adds ID, CreatedAt, UpdatedAt, DeletedAt
	ProductID                    uint    `gorm:"column:product_id" json:"product_id"`
	ProductItemName              string  `gorm:"column:product_item_name" json:"product_item_name"`
	MinimumLoanPeriod            int     `gorm:"column:minimum_loan_period" json:"minimum_loan_period"`
	MaximumLoanPeriod            int     `gorm:"column:maximum_loan_period" json:"maximum_loan_period"`
	MinimumLoanAmount            float64 `gorm:"column:minimum_loan_amount" json:"minimum_loan_amount"`
	MaximumLoanAmount            float64 `gorm:"column:maximum_loan_amount" json:"maximum_loan_amount"`
	InterestApplyType            string  `gorm:"column:interest_apply_type" json:"interest_apply_type"`
	MinimumInterest              float64 `gorm:"column:minimum_interest" json:"minimum_interest"`
	MaximumInterest              float64 `gorm:"column:maximum_interest" json:"maximum_interest"`
	MinimumCollectionPeriod      int     `gorm:"column:minimum_collection_period" json:"minimum_collection_period"`
	MaximumCollectionPeriod      int     `gorm:"column:maximum_collection_period" json:"maximum_collection_period"`
	PenaltyMethod                string  `gorm:"column:penalty_method" json:"penalty_method"`
	PenaltyApplyType             string  `gorm:"column:penalty_apply_type" json:"penalty_apply_type"`
	PenaltyPercentage            float64 `gorm:"column:penalty_percentage" json:"penalty_percentage"`
	PenaltyStartAfterDays        int     `gorm:"column:penalty_start_after_days" json:"penalty_start_after_days"`
	SavingAmount                 float64 `gorm:"column:saving_amount" json:"saving_amount"`
	SavingAmountType             string  `gorm:"column:saving_amount_type" json:"saving_amount_type"`
	SavingPayment                string  `gorm:"column:saving_payment" json:"saving_payment"`
	SavingAccountMonthlyInterest float64 `gorm:"column:saving_account_monthly_interest" json:"saving_account_monthly_interest"`
	SavingInterestCalType        string  `gorm:"column:saving_interest_cal_type" json:"saving_interest_cal_type"`
	RequiredGuaranteeCount       int     `gorm:"column:required_guarantee_count" json:"required_guarantee_count"`
	SavingInterestRate           float64 `gorm:"column:saving_interest_rate" json:"saving_interest_rate"`
	IsDifferentCollectionPeriod  bool    `gorm:"column:is_different_collection_period" json:"is_different_collection_period"`
}

// TableName explicitly overrides pluralization to match exactly
func (ProductHasItem) TableName() string {
	return "product_has_items"
}
