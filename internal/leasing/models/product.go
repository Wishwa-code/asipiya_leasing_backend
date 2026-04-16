package models

import "gorm.io/gorm"

// Product maps to the legacy Laravel `products` table
type Product struct {
	gorm.Model
	ProductName           string `gorm:"column:product_name" json:"product_name"`
	ProductCode           string `gorm:"column:product_code" json:"product_code"`
	InterestMethod        string `gorm:"column:interest_method" json:"interest_method"`
	LoanPeriodType        string `gorm:"column:loan_period_type" json:"loan_period_type"`
	InterestPeriodType    string `gorm:"column:interest_period_type" json:"interest_period_type"`
	CollectionPeriodType  string `gorm:"column:collection_period_type" json:"collection_period_type"`
	CollectionDateType    string `gorm:"column:collection_date_type" json:"collection_date_type"`
	GuaranteeCount        int    `gorm:"column:guarantee_count" json:"guarantee_count"`
	SavingAmountType      string `gorm:"column:saving_amount_type" json:"saving_amount_type"`
	SavingCollectionType  string `gorm:"column:saving_collection_type" json:"saving_collection_type"`
	SavingInterestCalType string `gorm:"column:saving_interest_cal_type" json:"saving_interest_cal_type"`
	SavingAccountStatus   string `gorm:"column:saving_account_status" json:"saving_account_status"`
	RecoveryAccountStatus string `gorm:"column:recovery_account_status" json:"recovery_account_status"`
	Status                string `gorm:"column:status" json:"status"`
	CreatedBy             *uint  `gorm:"column:created_by" json:"created_by"`
	UpdatedBy             *uint  `gorm:"column:updated_by" json:"updated_by"`
	DeletedBy             *uint  `gorm:"column:deleted_by" json:"deleted_by"`

	// Relationships
	ProductHasItems   []ProductHasItem           `gorm:"foreignKey:ProductID" json:"product_has_items"`
	AdditionalCharges []ProductAdditionalCharges `gorm:"foreignKey:ProductID" json:"additional_charges"`
	RequiredDocuments []ProductRequiredDocuments `gorm:"foreignKey:ProductID" json:"required_documents"`
}

// TableName explicitly overrides pluralization to match exactly
func (Product) TableName() string {
	return "products"
}
